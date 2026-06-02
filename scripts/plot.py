# /// script
# requires-python = ">=3.10"
# dependencies = [
#     "matplotlib",
# ]
# ///
"""Plot benchmark results produced by scripts/benchmark.py.

Reads `.benchmark_results.json` from the repo root and writes a boxplot
comparing execution time across directory sizes and modes.

Usage:
    scripts/plot.py [--out OUT.png]
"""

import argparse
import json
import sys
from pathlib import Path

import matplotlib.pyplot as plt

REPO_ROOT = Path(__file__).resolve().parent.parent
RESULTS_PATH = REPO_ROOT / ".benchmark_results.json"


def die(msg: str) -> None:
    print(f"error: {msg}", file=sys.stderr)
    sys.exit(1)


def plot_results(
    results: list[dict],
    series: list[tuple[str, str, str]],
    out_path: Path,
) -> None:
    plt.style.use("seaborn-v0_8-whitegrid")

    sizes = sorted({r["size"] for r in results})
    palette = ["#4C72B0", "#55A868", "#C44E52", "#8172B2"]
    series_keys = [(binary, mode) for binary, mode, _ in series]
    colors = {key: palette[i % len(palette)] for i, key in enumerate(series_keys)}
    markers = {key: m for key, m in zip(series_keys, ["o", "s", "D", "^"])}

    fig, ax = plt.subplots(figsize=(9, 5.5))

    for key in series_keys:
        means, stddevs = [], []
        for size in sizes:
            match = next(
                r
                for r in results
                if r["size"] == size and r["binary"] == key[0] and r["mode"] == key[1]
            )
            means.append(match["mean"])
            stddevs.append(match["stddev"])

        color = colors[key]
        ax.plot(
            sizes,
            means,
            marker=markers[key],
            markersize=7,
            linewidth=2,
            color=color,
            label=f"{key[0]} ({key[1]})",
            zorder=3,
        )
        lower = [m - s for m, s in zip(means, stddevs)]
        upper = [m + s for m, s in zip(means, stddevs)]
        ax.fill_between(sizes, lower, upper, color=color, alpha=0.15, zorder=2)

    ax.set_xlabel("Directory entries (target)")
    ax.set_ylabel("Wall time (ms)")
    ax.set_title("Execution time vs. directory size")
    ax.set_xticks(sizes)
    ax.set_ylim(bottom=0)
    ax.legend(loc="upper left", frameon=True)

    fig.tight_layout()
    fig.savefig(out_path, dpi=150)
    print(f">> plot saved to {out_path}")


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--results",
        default=str(RESULTS_PATH),
        help="path to benchmark results JSON (default: .benchmark_results.json in repo root)",
    )
    parser.add_argument(
        "--out",
        default=str(REPO_ROOT / "benchmark.png"),
        help="output path for the plot image (default: benchmark.png in repo root)",
    )
    args = parser.parse_args()

    results_path = Path(args.results).resolve()
    if not results_path.exists():
        die(f"{results_path} not found. Run scripts/benchmark.py first")

    payload = json.loads(results_path.read_text())
    series = [tuple(s) for s in payload["series"]]
    plot_results(payload["results"], series, Path(args.out).resolve())


if __name__ == "__main__":
    main()
