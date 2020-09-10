package main

var Icon_Ext = map[string]*iInfo{
	"htm":               iSet["html"],
	"html":              iSet["html"],
	"xhtml":             iSet["html"],
	"html_vm":           iSet["html"],
	"asp":               iSet["html"],
	"jade":              iSet["pug"],
	"pug":               iSet["pug"],
	"md":                iSet["markdown"],
	"markdown":          iSet["markdown"],
	"rst":               iSet["markdown"],
	"blink":             iSet["blink"],
	"css":               iSet["css"],
	"scss":              iSet["sass"],
	"sass":              iSet["sass"],
	"less":              iSet["less"],
	"json":              iSet["json"],
	"tsbuildinfo":       iSet["json"],
	"json5":             iSet["json"],
	"jsonl":             iSet["json"],
	"ndjson":            iSet["json"],
	"jinja":             iSet["jinja"],
	"jinja2":            iSet["jinja"],
	"j2":                iSet["jinja"],
	"jinja-html":        iSet["jinja"],
	"sublime-project":   iSet["sublime"],
	"sublime-workspace": iSet["sublime"],
	"yaml":              iSet["yaml"],
	"yaml-tmlanguage":   iSet["yaml"],
	"yml":               iSet["yaml"],
	"xml":               iSet["xml"],
	"plist":             iSet["xml"],
	"xsd":               iSet["xml"],
	"dtd":               iSet["xml"],
	"xsl":               iSet["xml"],
	"xslt":              iSet["xml"],
	"resx":              iSet["xml"],
	"iml":               iSet["xml"],
	"xquery":            iSet["xml"],
	"tmLanguage":        iSet["xml"],
	"manifest":          iSet["xml"],
	"project":           iSet["xml"],
	"png":               iSet["image"],
	"jpeg":              iSet["image"],
	"jpg":               iSet["image"],
	"gif":               iSet["image"],
	"ico":               iSet["image"],
	"tif":               iSet["image"],
	"tiff":              iSet["image"],
	"psd":               iSet["image"],
	"psb":               iSet["image"],
	"ami":               iSet["image"],
	"apx":               iSet["image"],
	"bmp":               iSet["image"],
	"bpg":               iSet["image"],
	"brk":               iSet["image"],
	"cur":               iSet["image"],
	"dds":               iSet["image"],
	"dng":               iSet["image"],
	"exr":               iSet["image"],
	"fpx":               iSet["image"],
	"gbr":               iSet["image"],
	"img":               iSet["image"],
	"jbig2":             iSet["image"],
	"jb2":               iSet["image"],
	"jng":               iSet["image"],
	"jxr":               iSet["image"],
	"pbm":               iSet["image"],
	"pgf":               iSet["image"],
	"pic":               iSet["image"],
	"raw":               iSet["image"],
	"webp":              iSet["image"],
	"eps":               iSet["image"],
	"afphoto":           iSet["image"],
	"ase":               iSet["image"],
	"aseprite":          iSet["image"],
	"clip":              iSet["image"],
	"cpt":               iSet["image"],
	"heif":              iSet["image"],
	"heic":              iSet["image"],
	"kra":               iSet["image"],
	"mdp":               iSet["image"],
	"ora":               iSet["image"],
	"pdn":               iSet["image"],
	"reb":               iSet["image"],
	"sai":               iSet["image"],
	"tga":               iSet["image"],
	"xcf":               iSet["image"],
	"js":                iSet["javascript"],
	"esx":               iSet["javascript"],
	"mj":                iSet["javascript"],
	"jsx":               iSet["react"],
	"tsx":               iSet["react_ts"],
	"ini":               iSet["settings"],
	"dlc":               iSet["settings"],
	"dll":               iSet["settings"],
	"config":            iSet["settings"],
	"conf":              iSet["settings"],
	"properties":        iSet["settings"],
	"prop":              iSet["settings"],
	"settings":          iSet["settings"],
	"option":            iSet["settings"],
	"props":             iSet["settings"],
	"toml":              iSet["settings"],
	"prefs":             iSet["settings"],
	"dotsettings":       iSet["settings"],
	"cfg":               iSet["settings"],
	"ts":                iSet["typescript"],
	"marko":             iSet["markojs"],
	"pdf":               iSet["pdf"],
	"xlsx":              iSet["table"],
	"xls":               iSet["table"],
	"csv":               iSet["table"],
	"tsv":               iSet["table"],
	"vscodeignore":      iSet["vscode"],
	"vsixmanifest":      iSet["vscode"],
	"vsix":              iSet["vscode"],
	"code-workplace":    iSet["vscode"],
	"csproj":            iSet["visualstudio"],
	"ruleset":           iSet["visualstudio"],
	"sln":               iSet["visualstudio"],
	"suo":               iSet["visualstudio"],
	"vb":                iSet["visualstudio"],
	"vbs":               iSet["visualstudio"],
	"vcxitems":          iSet["visualstudio"],
	"vcxproj":           iSet["visualstudio"],
	"pdb":               iSet["database"],
	"sql":               iSet["mysql"],
	"pks":               iSet["database"],
	"pkb":               iSet["database"],
	"accdb":             iSet["database"],
	"mdb":               iSet["database"],
	"sqlite":            iSet["sqlite"],
	"sqlite3":           iSet["sqlite"],
	"pgsql":             iSet["postgresql"],
	"postgres":          iSet["postgresql"],
	"psql":              iSet["postgresql"],
	"cs":                iSet["csharp"],
	"csx":               iSet["csharp"],
	"qs":                iSet["qsharp"],
	"zip":               iSet["zip"],
	"tar":               iSet["zip"],
	"gz":                iSet["zip"],
	"xz":                iSet["zip"],
	"br":                iSet["zip"],
	"bzip2":             iSet["zip"],
	"gzip":              iSet["zip"],
	"brotli":            iSet["zip"],
	"7z":                iSet["zip"],
	"rar":               iSet["zip"],
	"tgz":               iSet["zip"],
	"vala":              iSet["vala"],
	"zig":               iSet["zig"],
	"exe":               iSet["exe"],
	"msi":               iSet["exe"],
	"java":              iSet["java"],
	"jar":               iSet["java"],
	"jsp":               iSet["java"],
	"c":                 iSet["c"],
	"m":                 iSet["c"],
	"i":                 iSet["c"],
	"mi":                iSet["c"],
	"h":                 iSet["h"],
	"cc":                iSet["cpp"],
	"cpp":               iSet["cpp"],
	"cxx":               iSet["cpp"],
	"c++":               iSet["cpp"],
	"cp":                iSet["cpp"],
	"mm":                iSet["cpp"],
	"mii":               iSet["cpp"],
	"ii":                iSet["cpp"],
	"hh":                iSet["hpp"],
	"hpp":               iSet["hpp"],
	"hxx":               iSet["hpp"],
	"h++":               iSet["hpp"],
	"hp":                iSet["hpp"],
	"tcc":               iSet["hpp"],
	"inl":               iSet["hpp"],
	"go":                iSet["go"],
	"py":                iSet["python"],
	"pyc":               iSet["python-misc"],
	"whl":               iSet["python-misc"],
	"url":               iSet["url"],
	"sh":                iSet["console"],
	"ksh":               iSet["console"],
	"csh":               iSet["console"],
	"tcsh":              iSet["console"],
	"zsh":               iSet["console"],
	"bash":              iSet["console"],
	"bat":               iSet["console"],
	"cmd":               iSet["console"],
	"awk":               iSet["console"],
	"fish":              iSet["console"],
	"ps1":               iSet["powershell"],
	"psm1":              iSet["powershell"],
	"psd1":              iSet["powershell"],
	"ps1xml":            iSet["powershell"],
	"psc1":              iSet["powershell"],
	"pssc":              iSet["powershell"],
	"gradle":            iSet["gradle"],
	"doc":               iSet["word"],
	"docx":              iSet["word"],
	"rtf":               iSet["word"],
	"cer":               iSet["certificate"],
	"cert":              iSet["certificate"],
	"crt":               iSet["certificate"],
	"pub":               iSet["key"],
	"key":               iSet["key"],
	"pem":               iSet["key"],
	"asc":               iSet["key"],
	"gpg":               iSet["key"],
	"woff":              iSet["font"],
	"woff2":             iSet["font"],
	"ttf":               iSet["font"],
	"eot":               iSet["font"],
	"suit":              iSet["font"],
	"otf":               iSet["font"],
	"bmap":              iSet["font"],
	"fnt":               iSet["font"],
	"odttf":             iSet["font"],
	"ttc":               iSet["font"],
	"font":              iSet["font"],
	"fonts":             iSet["font"],
	"sui":               iSet["font"],
	"ntf":               iSet["font"],
	"mrf":               iSet["font"],
	"lib":               iSet["lib"],
	"bib":               iSet["lib"],
	"rb":                iSet["ruby"],
	"erb":               iSet["ruby"],
	"fs":                iSet["fsharp"],
	"fsx":               iSet["fsharp"],
	"fsi":               iSet["fsharp"],
	"fsproj":            iSet["fsharp"],
	"swift":             iSet["swift"],
	"ino":               iSet["arduino"],
	"dockerignore":      iSet["docker"],
	"dockerfile":        iSet["docker"],
	"tex":               iSet["tex"],
	"sty":               iSet["tex"],
	"dtx":               iSet["tex"],
	"ltx":               iSet["tex"],
	"pptx":              iSet["powerpoint"],
	"ppt":               iSet["powerpoint"],
	"pptm":              iSet["powerpoint"],
	"potx":              iSet["powerpoint"],
	"potm":              iSet["powerpoint"],
	"ppsx":              iSet["powerpoint"],
	"ppsm":              iSet["powerpoint"],
	"pps":               iSet["powerpoint"],
	"ppam":              iSet["powerpoint"],
	"ppa":               iSet["powerpoint"],
	"webm":              iSet["video"],
	"mkv":               iSet["video"],
	"flv":               iSet["video"],
	"vob":               iSet["video"],
	"ogv":               iSet["video"],
	"ogg":               iSet["video"],
	"gifv":              iSet["video"],
	"avi":               iSet["video"],
	"mov":               iSet["video"],
	"qt":                iSet["video"],
	"wmv":               iSet["video"],
	"yuv":               iSet["video"],
	"rm":                iSet["video"],
	"rmvb":              iSet["video"],
	"mp4":               iSet["video"],
	"m4v":               iSet["video"],
	"mpg":               iSet["video"],
	"mp2":               iSet["video"],
	"mpeg":              iSet["video"],
	"mpe":               iSet["video"],
	"mpv":               iSet["video"],
	"m2v":               iSet["video"],
	"vdi":               iSet["virtual"],
	"vbox":              iSet["virtual"],
	"vbox-prev":         iSet["virtual"],
	"ics":               iSet["email"],
	"mp3":               iSet["audio"],
	"flac":              iSet["audio"],
	"m4a":               iSet["audio"],
	"wma":               iSet["audio"],
	"aiff":              iSet["audio"],
	"coffee":            iSet["coffee"],
	"cson":              iSet["coffee"],
	"iced":              iSet["coffee"],
	"txt":               iSet["document"],
	"graphql":           iSet["graphql"],
	"gql":               iSet["graphql"],
	"rs":                iSet["rust"],
	"raml":              iSet["raml"],
	"xaml":              iSet["xaml"],
	"hs":                iSet["haskell"],
	"kt":                iSet["kotlin"],
	"kts":               iSet["kotlin"],
	"patch":             iSet["git"],
	"lua":               iSet["lua"],
	"clj":               iSet["clojure"],
	"cljs":              iSet["clojure"],
	"cljc":              iSet["clojure"],
	"groovy":            iSet["groovy"],
	"r":                 iSet["r"],
	"rmd":               iSet["r"],
	"dart":              iSet["dart"],
	"as":                iSet["actionscript"],
	"mxml":              iSet["mxml"],
	"ahk":               iSet["autohotkey"],
	"swf":               iSet["flash"],
	"swc":               iSet["swc"],
	"cmake":             iSet["cmake"],
	"asm":               iSet["assembly"],
	"a51":               iSet["assembly"],
	"inc":               iSet["assembly"],
	"nasm":              iSet["assembly"],
	"s":                 iSet["assembly"],
	"ms":                iSet["assembly"],
	"agc":               iSet["assembly"],
	"ags":               iSet["assembly"],
	"aea":               iSet["assembly"],
	"argus":             iSet["assembly"],
	"mitigus":           iSet["assembly"],
	"binsource":         iSet["assembly"],
	"vue":               iSet["vue"],
	"ml":                iSet["ocaml"],
	"mli":               iSet["ocaml"],
	"cmx":               iSet["ocaml"],
	"lock":              iSet["lock"],
	"hbs":               iSet["handlebars"],
	"mustache":          iSet["handlebars"],
	"pm":                iSet["perl"],
	"raku":              iSet["perl"],
	"hx":                iSet["haxe"],
	"pp":                iSet["puppet"],
	"ex":                iSet["elixir"],
	"exs":               iSet["elixir"],
	"eex":               iSet["elixir"],
	"leex":              iSet["elixir"],
	"erl":               iSet["erlang"],
	"twig":              iSet["twig"],
	"jl":                iSet["julia"],
	"elm":               iSet["elm"],
	"pure":              iSet["purescript"],
	"purs":              iSet["purescript"],
	"tpl":               iSet["smarty"],
	"styl":              iSet["stylus"],
	"merlin":            iSet["merlin"],
	"v":                 iSet["verilog"],
	"vhd":               iSet["verilog"],
	"sv":                iSet["verilog"],
	"svh":               iSet["verilog"],
	"robot":             iSet["robot"],
	"sol":               iSet["solidity"],
	"yang":              iSet["yang"],
	"mjml":              iSet["mjml"],
	"tf":                iSet["terraform"],
	"tfvars":            iSet["terraform"],
	"tfstate":           iSet["terraform"],
	"applescript":       iSet["applescript"],
	"ipa":               iSet["applescript"],
	"cake":              iSet["cake"],
	"nim":               iSet["nim"],
	"nimble":            iSet["nim"],
	"apib":              iSet["apiblueprint"],
	"apiblueprint":      iSet["apiblueprint"],
	"pcss":              iSet["postcss"],
	"sss":               iSet["postcss"],
	"todo":              iSet["todo"],
	"nix":               iSet["nix"],
	"slim":              iSet["slim"],
	"http":              iSet["http"],
	"rest":              iSet["http"],
	"apk":               iSet["android"],
	"env":               iSet["tune"],
	"jenkinsfile":       iSet["jenkins"],
	"jenkins":           iSet["jenkins"],
	"log":               iSet["log"],
	"ejs":               iSet["ejs"],
	"djt":               iSet["django"],
	"pot":               iSet["i18n"],
	"po":                iSet["i18n"],
	"mo":                iSet["i18n"],
	"d":                 iSet["d"],
	"mdx":               iSet["mdx"],
	"gd":                iSet["godot"],
	"godot":             iSet["godot-assets"],
	"tres":              iSet["godot-assets"],
	"tscn":              iSet["godot-assets"],
	"azcli":             iSet["azure"],
	"vagrantfile":       iSet["vagrant"],
	"cshtml":            iSet["razor"],
	"vbhtml":            iSet["razor"],
	"ad":                iSet["asciidoc"],
	"adoc":              iSet["asciidoc"],
	"asciidoc":          iSet["asciidoc"],
	"edge":              iSet["edge"],
	"ss":                iSet["scheme"],
	"scm":               iSet["scheme"],
	"stl":               iSet["3d"],
	"obj":               iSet["3d"],
	"ac":                iSet["3d"],
	"blend":             iSet["3d"],
	"mesh":              iSet["3d"],
	"mqo":               iSet["3d"],
	"pmd":               iSet["3d"],
	"pmx":               iSet["3d"],
	"skp":               iSet["3d"],
	"vac":               iSet["3d"],
	"vdp":               iSet["3d"],
	"vox":               iSet["3d"],
	"svg":               iSet["svg"],
	"vimrc":             iSet["vim"],
	"gvimrc":            iSet["vim"],
	"exrc":              iSet["vim"],
	"moon":              iSet["moonscript"],
	"iso":               iSet["disc"],
	"f":                 iSet["fortran"],
	"f77":               iSet["fortran"],
	"f90":               iSet["fortran"],
	"f95":               iSet["fortran"],
	"f03":               iSet["fortran"],
	"f08":               iSet["fortran"],
	"tcl":               iSet["tcl"],
	"liquid":            iSet["liquid"],
	"p":                 iSet["prolog"],
	"pro":               iSet["prolog"],
	"coco":              iSet["coconut"],
	"sketch":            iSet["sketch"],
	"opam":              iSet["opam"],
	"dhallb":            iSet["dhall"],
	"pwn":               iSet["pawn"],
	"amx":               iSet["pawn"],
	"dhall":             iSet["dhall"],
	"pas":               iSet["pascal"],
	"unity":             iSet["shaderlab"],
	"nupkg":             iSet["nuget"],
	"command":           iSet["command"],
	"dsc":               iSet["denizenscript"],
	"deb":               iSet["debian"],
	"rpm":               iSet["redhat"],
	"snap":              iSet["ubuntu"],
	"ebuild":            iSet["gentoo"],
	"pkg":               iSet["applescript"],
	"openbsd":           iSet["freebsd"],
	// "ls":                iSet["livescript"],
	// "re":                iSet["reason"],
	// "rei":               iSet["reason"],
	// "cmj":               iSet["bucklescript"],
	// "nb":                iSet["mathematica"],
	// "wl":                iSet["wolframlanguage"],
	// "wls":               iSet["wolframlanguage"],
	// "njk":               iSet["nunjucks"],
	// "nunjucks":          iSet["nunjucks"],
	// "au3":               iSet["autoit"],
	// "haml":              iSet["haml"],
	// "feature":           iSet["cucumber"],
	// "riot":              iSet["riot"],
	// "tag":               iSet["riot"],
	// "vfl":               iSet["vfl"],
	// "kl":                iSet["kl"],
	// "cfml":              iSet["coldfusion"],
	// "cfc":               iSet["coldfusion"],
	// "lucee":             iSet["coldfusion"],
	// "cfm":               iSet["coldfusion"],
	// "cabal":             iSet["cabal"],
	// "rql":               iSet["restql"],
	// "restql":            iSet["restql"],
	// "kv":                iSet["kivy"],
	// "graphcool":         iSet["graphcool"],
	// "sbt":               iSet["sbt"],
	// "cr":                iSet["crystal"],
	// "ecr":               iSet["crystal"],
	// "cu":                iSet["cuda"],
	// "cuh":               iSet["cuda"],
	// "def":               iSet["dotjs"],
	// "dot":               iSet["dotjs"],
	// "jst":               iSet["dotjs"],
	// "pde":               iSet["processing"],
	// "wpy":               iSet["wepy"],
	// "hcl":               iSet["hcl"],
	// "san":               iSet["san"],
	// "red":               iSet["red"],
	// "fxp":               iSet["foxpro"],
	// "prg":               iSet["foxpro"],
	// "wat":               iSet["webassembly"],
	// "wasm":              iSet["webassembly"],
	// "ipynb":             iSet["jupyter"],
	// "bal":               iSet["ballerina"],
	// "balx":              iSet["ballerina"],
	// "rkt":               iSet["racket"],
	// "bzl":               iSet["bazel"],
	// "bazel":             iSet["bazel"],
	// "mint":              iSet["mint"],
	// "vm":                iSet["velocity"],
	// "fhtml":             iSet["velocity"],
	// "vtl":               iSet["velocity"],
	// "prisma":            iSet["prisma"],
	// "abc":               iSet["abc"],
	// "lisp":              iSet["lisp"],
	// "lsp":               iSet["lisp"],
	// "cl":                iSet["lisp"],
	// "fast":              iSet["lisp"],
	// "svelte":            iSet["svelte"],
	// "prw":               iSet["advpl_prw"],
	// "prx":               iSet["advpl_prw"],
	// "ptm":               iSet["advpl_ptm"],
	// "tlpp":              iSet["advpl_tlpp"],
	// "ch":                iSet["advpl_include"],
	// "4th":               iSet["forth"],
	// "fth":               iSet["forth"],
	// "frt":               iSet["forth"],
	// "iuml":              iSet["uml"],
	// "pu":                iSet["uml"],
	// "puml":              iSet["uml"],
	// "plantuml":          iSet["uml"],
	// "wsd":               iSet["uml"],
	// "sml":               iSet["sml"],
	// "mlton":             iSet["sml"],
	// "mlb":               iSet["sml"],
	// "sig":               iSet["sml"],
	// "fun":               iSet["sml"],
	// "cm":                iSet["sml"],
	// "lex":               iSet["sml"],
	// "use":               iSet["sml"],
	// "grm":               iSet["sml"],
	// "imba":              iSet["imba"],
	// "drawio":            iSet["drawio"],
	// "dio":               iSet["drawio"],
	// "sas":               iSet["sas"],
	// "sas7bdat":          iSet["sas"],
	// "sashdat":           iSet["sas"],
	// "astore":            iSet["sas"],
	// "ast":               iSet["sas"],
	// "sast":              iSet["sas"],
}
