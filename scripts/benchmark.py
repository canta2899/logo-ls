# /// script
# requires-python = ">=3.10"
# dependencies = []
# ///
"""Benchmark a logo-ls binary across multiple directory sizes.

For each size, the binary is benchmarked with and without `-D` (git status),
and compared against `ls` (plain mode only, since ls has no git integration).
Prints a text table and writes raw results to `.benchmark_results.json` in
the repo root for later plotting.

Usage:
    scripts/benchmark.py [--logo-ls PATH] [--ls PATH]
"""

import argparse
import json
import os
import shutil
import stat
import subprocess
import sys
import tempfile
from pathlib import Path

RUNS = 300
WARMUP = 5
SIZES = [100, 500, 1500, 3000]
REPO_ROOT = Path(__file__).resolve().parent.parent
RESULTS_PATH = REPO_ROOT / ".benchmark_results.json"


def die(msg: str) -> None:
    print(f"error: {msg}", file=sys.stderr)
    sys.exit(1)


def touch(path: Path, executable: bool = False) -> None:
    path.touch()
    if executable:
        path.chmod(path.stat().st_mode | stat.S_IXUSR | stat.S_IXGRP | stat.S_IXOTH)


def build_test_dir(bench_dir: Path, target_entries: int) -> int:
    """Seed `bench_dir` with roughly `target_entries` filesystem entries."""
    for child in bench_dir.iterdir():
        if child.is_dir() and not child.is_symlink():
            shutil.rmtree(child)
        else:
            child.unlink()

    subprocess.run(["git", "init", "-q"], cwd=bench_dir, check=True)
    subprocess.run(
        ["git", "config", "user.email", "bench@test.com"], cwd=bench_dir, check=True
    )
    subprocess.run(["git", "config", "user.name", "Bench"], cwd=bench_dir, check=True)

    n_files = int(target_entries * 0.60)
    n_dotfiles = int(target_entries * 0.20)
    n_dirs = max(1, int(target_entries * 0.10))
    n_links = int(target_entries * 0.05)
    n_execs = int(target_entries * 0.05)

    exts = ["txt", "json", "png", "sh", "go", "md"]
    for i in range(n_files):
        ext = exts[i % len(exts)]
        touch(bench_dir / f"file_{i}.{ext}", executable=(ext == "sh"))

    for i in range(n_dotfiles):
        touch(bench_dir / f".dotfile_{i}")

    for i in range(n_dirs):
        d = bench_dir / f"dir_{i}"
        d.mkdir()
        for j in range(5):
            touch(d / f"nested_{j}.txt")

    for i in range(n_links):
        target = f"file_{i}.txt" if i < n_files else f"missing_{i}"
        (bench_dir / f"link_{i}").symlink_to(target)

    for i in range(n_execs):
        touch(bench_dir / f"bin_{i}", executable=True)

    to_add = [
        p.name
        for p in list(bench_dir.iterdir())[: target_entries // 2]
        if not p.name.startswith(".git")
    ]
    if to_add:
        subprocess.run(["git", "add", *to_add], cwd=bench_dir, check=True)
        subprocess.run(
            ["git", "commit", "-q", "-m", "initial"], cwd=bench_dir, check=True
        )
        modified = [p for p in to_add if (bench_dir / p).is_file()][
            : max(1, target_entries // 20)
        ]
        for name in modified:
            with (bench_dir / name).open("a") as f:
                f.write("modified\n")

    return sum(1 for _ in bench_dir.iterdir())


def run_hyperfine(binary: str, bench_dir: Path, flags: str) -> dict:
    with tempfile.NamedTemporaryFile(suffix=".json", delete=False) as tmp:
        out_path = tmp.name
    try:
        cmd = [
            "hyperfine",
            "--warmup",
            str(WARMUP),
            "--shell=none",
            "--runs",
            str(RUNS),
            "--export-json",
            out_path,
            f"{binary} {flags} {bench_dir}".strip(),
        ]
        subprocess.run(
            cmd, check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL
        )
        data = json.load(open(out_path))
    finally:
        os.unlink(out_path)

    r = data["results"][0]
    return {
        "mean": r["mean"] * 1000,
        "stddev": r["stddev"] * 1000,
        "min": r["min"] * 1000,
        "max": r["max"] * 1000,
        "times": [t * 1000 for t in r["times"]],
    }


def print_table(results: list[dict]) -> None:
    headers = (
        "Binary",
        "Mode",
        "Size",
        "Mean (ms)",
        "Stddev (ms)",
        "Min (ms)",
        "Max (ms)",
    )
    rows: list[tuple[str, ...]] = [headers]
    for r in results:
        rows.append(
            (
                r["binary"],
                r["mode"],
                str(r["size"]),
                f"{r['mean']:.2f}",
                f"{r['stddev']:.2f}",
                f"{r['min']:.2f}",
                f"{r['max']:.2f}",
            )
        )

    widths = [max(len(row[i]) for row in rows) for i in range(len(headers))]
    sep = "-+-".join("-" * w for w in widths)

    def fmt(row):
        return " | ".join(c.ljust(widths[i]) for i, c in enumerate(row))

    print(fmt(rows[0]))
    print(sep)
    for row in rows[1:]:
        print(fmt(row))


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--logo-ls",
        default=str(REPO_ROOT / "logo-ls"),
        help="path to the logo-ls binary (default: ./logo-ls in repo root)",
    )
    parser.add_argument(
        "--ls",
        default="/bin/ls",
        help="path to the ls binary to compare against (default: /bin/ls)",
    )
    args = parser.parse_args()

    binary = Path(args.logo_ls).resolve()
    if not binary.is_file() or not os.access(binary, os.X_OK):
        die(f"logo-ls binary '{binary}' not found or not executable")

    ls_path = Path(args.ls).resolve()
    if not ls_path.is_file() or not os.access(ls_path, os.X_OK):
        die(f"ls binary '{ls_path}' not found or not executable")

    if not shutil.which("hyperfine"):
        die("hyperfine is required (brew install hyperfine)")

    series: list[tuple[str, str, str]] = [
        ("ls", "plain", ""),
        ("logo-ls", "plain", ""),
        ("logo-ls", "git", "-D"),
    ]
    bin_for_label = {"ls": str(ls_path), "logo-ls": str(binary)}

    print(f">> logo-ls: {binary}")
    print(f">> ls:      {ls_path}")
    print(f">> sizes:   {SIZES}")
    print(f">> {RUNS} runs per configuration, {WARMUP} warmup")

    bench_dir = Path(tempfile.mkdtemp(prefix="logo-ls-bench-"))
    print(f">> bench dir: {bench_dir}")

    results = []
    try:
        for size in SIZES:
            actual = build_test_dir(bench_dir, size)
            print(f">> seeded {actual} entries")
            for label, mode, flags in series:
                print(f"   size={size:5d} binary={label} mode={mode}...")
                r = run_hyperfine(bin_for_label[label], bench_dir, flags)
                r["size"] = size
                r["binary"] = label
                r["mode"] = mode
                results.append(r)
    finally:
        if bench_dir.exists():
            shutil.rmtree(bench_dir)
            print(f">> cleaned up {bench_dir}")

    print()
    print_table(results)

    payload = {
        "runs": RUNS,
        "warmup": WARMUP,
        "sizes": SIZES,
        "series": series,
        "results": results,
    }
    RESULTS_PATH.write_text(json.dumps(payload, indent=2))
    print(f"\n>> results written to {RESULTS_PATH}")


if __name__ == "__main__":
    main()
