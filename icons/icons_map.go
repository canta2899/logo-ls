// assets contains all the Icon glyphs info
package icons

import "fmt"

// Icon_Info (icon information)
type IconInfo struct {
	i string
	c [3]uint8 // represents the color in rgb (default 0,0,0 is black)
	e bool     // whether or not the file is executable [true = is executable]
}

func (i *IconInfo) GetGlyph() string {
	return i.i
}

func (i *IconInfo) GetColor(f uint8) string {
	if i.e {
		return "\033[38;2;76;175;080m"
	} else if f == 1 {
		return fmt.Sprintf("\033[38;2;%03d;%03d;%03dm", i.c[0], i.c[1], i.c[2])
	}
	return fmt.Sprintf("\033[38;2;%03d;%03d;%03dm", i.c[0], i.c[1], i.c[2])
}

func (i *IconInfo) MakeExe() {
	i.e = true
}

var IconSet = map[string]*IconInfo{
	"3d":               {i: "\U0000e79b", c: [3]uint8{40, 182, 246}},  // 3d
	"actionscript":     {i: "\U0000fb25", c: [3]uint8{244, 68, 62}},   // actionscript (Not supported by nerdFont)
	"alpine":           {i: "\U0000f300", c: [3]uint8{14, 87, 123}},   // alpine
	"android":          {i: "\U000f0032", c: [3]uint8{139, 195, 74}},  // android
	"apiblueprint":     {i: "\U0000f031", c: [3]uint8{66, 165, 245}},  // apiblueprint (Not supported by nerdFont)
	"applescript":      {i: "\U0000f302", c: [3]uint8{120, 144, 156}}, // applescript
	"arch":             {i: "\U0000f303", c: [3]uint8{33, 142, 202}},  // arch
	"arduino":          {i: "\U0000f34b", c: [3]uint8{35, 151, 156}},  // arduino
	"asciidoc":         {i: "\U000f0219", c: [3]uint8{244, 68, 62}},   // asciidoc
	"assembly":         {i: "\U0000f471", c: [3]uint8{250, 109, 63}},  // assembly
	"audio":            {i: "\U000f0388", c: [3]uint8{239, 83, 80}},   // audio
	"authors":          {i: "\U0000f0c0", c: [3]uint8{244, 68, 62}},   // authors
	"autohotkey":       {i: "\U000f0313", c: [3]uint8{76, 175, 80}},   // autohotkey (Not supported by nerdFont)
	"azure":            {i: "\U0000fd03", c: [3]uint8{31, 136, 229}},  // azure
	"azure-pipelines":  {i: "\U0000f427", c: [3]uint8{20, 101, 192}},  // azure-pipelines
	"babel":            {i: "\U000f0a25", c: [3]uint8{253, 217, 59}},  // babel
	"bitbucket":        {i: "\U0000f171", c: [3]uint8{31, 136, 229}},  // bitbucket
	"blink":            {i: "\U000f022b", c: [3]uint8{249, 169, 60}},  // blink (The Foundry Nuke) (Not supported by nerdFont)
	"bower":            {i: "\U0000e61a", c: [3]uint8{239, 88, 60}},   // bower
	"bun":              {i: "\U000f0685", c: [3]uint8{249, 241, 225}}, // bun
	"c":                {i: "\U000f0671", c: [3]uint8{2, 119, 189}},   // c
	"cake":             {i: "\U000f00eb", c: [3]uint8{250, 111, 66}},  // cake
	"centOS":           {i: "\U0000f304", c: [3]uint8{157, 83, 135}},  // centOS
	"certificate":      {i: "\U000f0124", c: [3]uint8{249, 89, 63}},   // certificate
	"changelog":        {i: "\U000f099b", c: [3]uint8{139, 195, 74}},  // changelog
	"clojure":          {i: "\U0000e76a", c: [3]uint8{100, 221, 23}},  // clojure
	"cmake":            {i: "\U0000f425", c: [3]uint8{178, 178, 179}}, // cmake (Not supported by nerdFont)
	"coconut":          {i: "\U000f00d3", c: [3]uint8{141, 110, 99}},  // coconut
	"code-climate":     {i: "\U000f0509", c: [3]uint8{238, 238, 238}}, // code-climate
	"codecov":          {i: "\U0000e37c", c: [3]uint8{237, 80, 122}},  // codecov (Not supported by nerdFont)
	"codeowners":       {i: "\U000f0008", c: [3]uint8{175, 180, 43}},  // codeowners
	"coffee":           {i: "\U000f0176", c: [3]uint8{66, 165, 245}},  // coffee
	"command":          {i: "\U000f0633", c: [3]uint8{175, 188, 194}}, // command
	"commitlint":       {i: "\U000f0718", c: [3]uint8{43, 150, 137}},  // commitlint
	"conduct":          {i: "\U000f014e", c: [3]uint8{205, 220, 57}},  // conduct
	"console":          {i: "\U000f018d", c: [3]uint8{250, 111, 66}},  // console
	"contributing":     {i: "\U000f014d", c: [3]uint8{255, 202, 61}},  // contributing
	"cpp":              {i: "\U000f0672", c: [3]uint8{2, 119, 189}},   // cpp
	"credits":          {i: "\U000f0260", c: [3]uint8{156, 204, 101}}, // credits
	"csharp":           {i: "\U000f031b", c: [3]uint8{2, 119, 189}},   // csharp
	"css":              {i: "\U000f031c", c: [3]uint8{66, 165, 245}},  // css
	"css-map":          {i: "\U0000e749", c: [3]uint8{66, 165, 245}},  // css-map
	"d":                {i: "\U0000e7af", c: [3]uint8{244, 68, 62}},   // d
	"dart":             {i: "\U0000e798", c: [3]uint8{87, 182, 240}},  // dart
	"database":         {i: "\U0000e706", c: [3]uint8{255, 202, 61}},  // database
	"debian":           {i: "\U0000f306", c: [3]uint8{211, 61, 76}},   // debian
	"denizenscript":    {i: "\U000f0af1", c: [3]uint8{255, 213, 79}},  // denizenscript (Not supported by nerdFont)
	"dhall":            {i: "\U0000f448", c: [3]uint8{120, 144, 156}}, // dhall
	"disc":             {i: "\U0000e271", c: [3]uint8{176, 190, 197}}, // disc
	"django":           {i: "\U0000e71d", c: [3]uint8{67, 160, 71}},   // django
	"docker":           {i: "\U0000f308", c: [3]uint8{1, 135, 201}},   // docker
	"document":         {i: "\U000f0219", c: [3]uint8{66, 165, 245}},  // document
	"dune":             {i: "\U000f0509", c: [3]uint8{244, 127, 61}},  // dune
	"edge":             {i: "\U000f0065", c: [3]uint8{239, 111, 60}},  // edge
	"ejs":              {i: "\U0000e618", c: [3]uint8{255, 202, 61}},  // ejs
	"elixir":           {i: "\U0000e62d", c: [3]uint8{149, 117, 205}}, // elixir
	"elm":              {i: "\U0000e62c", c: [3]uint8{96, 181, 204}},  // elm
	"email":            {i: "\U000f01ee", c: [3]uint8{66, 165, 245}},  // email
	"erlang":           {i: "\U0000e7b1", c: [3]uint8{244, 68, 62}},   // erlang
	"eslint":           {i: "\U000f0c7a", c: [3]uint8{121, 134, 203}}, // eslint
	"exe":              {i: "\U0000f2d0", c: [3]uint8{229, 77, 58}},   // exe
	"fastlane":         {i: "\U000f0700", c: [3]uint8{149, 119, 232}}, // fastlane (Not supported by nerdFont)
	"favicon":          {i: "\U0000e623", c: [3]uint8{255, 213, 79}},  // favicon
	"fedora":           {i: "\U0000f30a", c: [3]uint8{52, 103, 172}},  // fedora
	"firebase":         {i: "\U000f0967", c: [3]uint8{251, 193, 60}},  // firebase
	"flash":            {i: "\U000f0241", c: [3]uint8{198, 52, 54}},   // flash (Not supported by nerdFont)
	"font":             {i: "\U0000f031", c: [3]uint8{244, 68, 62}},   // font
	"fortran":          {i: "\U000f121a", c: [3]uint8{250, 111, 66}},  // fortran
	"freebsd":          {i: "\U0000f30c", c: [3]uint8{175, 44, 42}},   // freebsd
	"fsharp":           {i: "\U0000e7a7", c: [3]uint8{55, 139, 186}},  // fsharp
	"gatsby":           {i: "\U000f0e43", c: [3]uint8{70, 0, 130}},    // gatsby
	"gcp":              {i: "\U000f0163", c: [3]uint8{70, 136, 250}},  // gcp (Not supported by nerdFont)
	"gemfile":          {i: "\U0000e21e", c: [3]uint8{229, 61, 58}},   // gemfile
	"gentoo":           {i: "\U0000f30d", c: [3]uint8{148, 141, 211}}, // gentoo
	"git":              {i: "\U0000e702", c: [3]uint8{229, 77, 58}},   // git
	"gitlab":           {i: "\U0000f296", c: [3]uint8{226, 69, 57}},   // gitlab
	"go":               {i: "\U000f07d3", c: [3]uint8{32, 173, 194}},  // go
	"go-mod":           {i: "\U000f07d3", c: [3]uint8{237, 80, 122}},  // go-mod
	"go-test":          {i: "\U000f07d3", c: [3]uint8{255, 213, 79}},  // go-test
	"godot":            {i: "\U0000e65f", c: [3]uint8{79, 195, 247}},  // godot
	"godot-assets":     {i: "\U0000e65f", c: [3]uint8{129, 199, 132}}, // godot-assets
	"gradle":           {i: "\U0000e660", c: [3]uint8{29, 151, 167}},  // gradle
	"graphql":          {i: "\U000f0877", c: [3]uint8{237, 80, 122}},  // graphql
	"groovy":           {i: "\U0000f2a6", c: [3]uint8{41, 198, 218}},  // groovy
	"grunt":            {i: "\U0000e611", c: [3]uint8{251, 170, 61}},  // grunt
	"gulp":             {i: "\U0000e763", c: [3]uint8{229, 61, 58}},   // gulp
	"h":                {i: "\U000f0c00", c: [3]uint8{2, 119, 189}},   // h (Not supported by nerdFont)
	"handlebars":       {i: "\U0000e60f", c: [3]uint8{250, 111, 66}},  // handlebars
	"haskell":          {i: "\U0000e61f", c: [3]uint8{254, 168, 62}},  // haskell
	"haxe":             {i: "\U0000e666", c: [3]uint8{246, 137, 61}},  // haxe
	"helm":             {i: "\U000f0833", c: [3]uint8{32, 173, 194}},  // helm (Not supported by nerdFont)
	"heroku":           {i: "\U0000e607", c: [3]uint8{105, 99, 185}},  // heroku
	"hpp":              {i: "\U000f0b0f", c: [3]uint8{2, 119, 189}},   // hpp (Not supported by nerdFont)
	"html":             {i: "\U0000f13b", c: [3]uint8{228, 79, 57}},   // html
	"http":             {i: "\U0000f484", c: [3]uint8{66, 165, 245}},  // http
	"husky":            {i: "\U000f03e9", c: [3]uint8{229, 229, 229}}, // husky
	"i18n":             {i: "\U000f05ca", c: [3]uint8{121, 134, 203}}, // i18n (Not supported by nerdFont)
	"image":            {i: "\U000f021f", c: [3]uint8{48, 166, 154}},  // image
	"ionic":            {i: "\U0000e7a9", c: [3]uint8{79, 143, 247}},  // ionic
	"java":             {i: "\U000f0176", c: [3]uint8{244, 68, 62}},   // java
	"javascript":       {i: "\U0000e74e", c: [3]uint8{255, 202, 61}},  // javascript
	"javascript-map":   {i: "\U0000e781", c: [3]uint8{255, 202, 61}},  // javascript-map
	"jenkins":          {i: "\U0000e767", c: [3]uint8{240, 214, 183}}, // jenkins
	"jest":             {i: "\U000f0af7", c: [3]uint8{244, 85, 62}},   // jest (Not supported by nerdFont)
	"jinja":            {i: "\U0000e66f", c: [3]uint8{174, 44, 42}},   // jinja
	"json":             {i: "\U0000e60b", c: [3]uint8{251, 193, 60}},  // json
	"julia":            {i: "\U0000e624", c: [3]uint8{134, 82, 159}},  // julia
	"karma":            {i: "\U0000e622", c: [3]uint8{60, 190, 174}},  // karma
	"key":              {i: "\U000f0306", c: [3]uint8{48, 166, 154}},  // key
	"kotlin":           {i: "\U0000e634", c: [3]uint8{139, 195, 74}},  // kotlin
	"laravel":          {i: "\U0000e73f", c: [3]uint8{248, 80, 81}},   // laravel
	"less":             {i: "\U0000e60b", c: [3]uint8{2, 119, 189}},   // less
	"lib":              {i: "\U0000f831", c: [3]uint8{139, 195, 74}},  // lib
	"linux":            {i: "\U0000e712", c: [3]uint8{238, 207, 55}},  // linux
	"liquid":           {i: "\U0000e275", c: [3]uint8{40, 182, 246}},  // liquid
	"lock":             {i: "\U000f033e", c: [3]uint8{255, 213, 79}},  // lock
	"log":              {i: "\U0000f719", c: [3]uint8{175, 180, 43}},  // log
	"lua":              {i: "\U0000e620", c: [3]uint8{66, 165, 245}},  // lua
	"makefile":         {i: "\U000f0229", c: [3]uint8{239, 83, 80}},   // makefile
	"manjaro":          {i: "\U0000f312", c: [3]uint8{73, 185, 90}},   // manjaro
	"markdown":         {i: "\U0000eb1d", c: [3]uint8{66, 165, 245}},  // markdown
	"markojs":          {i: "\U0000f13b", c: [3]uint8{2, 119, 189}},   // markojs (Not supported by nerdFont)
	"mdx":              {i: "\U0000f853", c: [3]uint8{255, 202, 61}},  // mdx
	"merlin":           {i: "\U0000f136", c: [3]uint8{66, 165, 245}},  // merlin (Not supported by nerdFont)
	"midi":             {i: "\U000f08f2", c: [3]uint8{66, 165, 245}},  // midi
	"mint":             {i: "\U0000f30f", c: [3]uint8{125, 190, 58}},  // mint
	"mjml":             {i: "\U0000e714", c: [3]uint8{249, 89, 63}},   // mjml (Not supported by nerdFont)
	"mocha":            {i: "\U000f01aa", c: [3]uint8{161, 136, 127}}, // mocha (Not supported by nerdFont)
	"modernizr":        {i: "\U0000e720", c: [3]uint8{234, 72, 99}},   // modernizr
	"moonscript":       {i: "\U0000f186", c: [3]uint8{251, 193, 60}},  // moonscript
	"mxml":             {i: "\U0000f72d", c: [3]uint8{254, 168, 62}},  // mxml
	"mysql":            {i: "\U0000e704", c: [3]uint8{1, 94, 134}},    // mysql
	"nim":              {i: "\U000f01a5", c: [3]uint8{255, 202, 61}},  // nim
	"nix":              {i: "\U0000f313", c: [3]uint8{80, 117, 193}},  // nix
	"nodejs":           {i: "\U000f0399", c: [3]uint8{139, 195, 74}},  // nodejs
	"npm":              {i: "\U0000e71e", c: [3]uint8{203, 56, 55}},   // npm
	"nuget":            {i: "\U0000e77f", c: [3]uint8{3, 136, 209}},   // nuget (Not supported by nerdFont)
	"nuxt":             {i: "\U000f1106", c: [3]uint8{65, 184, 131}},  // nuxt
	"ocaml":            {i: "\U0000e67a", c: [3]uint8{253, 154, 62}},  // ocaml
	"opam":             {i: "\U0000f1ce", c: [3]uint8{255, 213, 79}},  // opam (Not supported by nerdFont)
	"opensuse":         {i: "\U0000f314", c: [3]uint8{111, 180, 36}},  // opensuse
	"pascal":           {i: "\U000f03db", c: [3]uint8{3, 136, 209}},   // pascal (Not supported by nerdFont)
	"pawn":             {i: "\U0000e261", c: [3]uint8{239, 111, 60}},  // pawn
	"pdf":              {i: "\U0000e67d", c: [3]uint8{244, 68, 62}},   // pdf
	"perl":             {i: "\U0000e769", c: [3]uint8{149, 117, 205}}, // perl
	"php":              {i: "\U0000e608", c: [3]uint8{65, 129, 190}},  // php
	"postcss":          {i: "\U000f031c", c: [3]uint8{244, 68, 62}},   // postcss (Not supported by nerdFont)
	"postgresql":       {i: "\U0000e76e", c: [3]uint8{49, 99, 140}},   // postgresql
	"powerpoint":       {i: "\U000f0227", c: [3]uint8{209, 71, 51}},   // powerpoint
	"powershell":       {i: "\U000f0a0a", c: [3]uint8{5, 169, 244}},   // powershell
	"prettier":         {i: "\U0000e6b4", c: [3]uint8{86, 179, 180}},  // prettier
	"prisma":           {i: "\U0000e684", c: [3]uint8{229, 56, 150}},  // prisma
	"prolog":           {i: "\U0000e7a1", c: [3]uint8{239, 83, 80}},   // prolog
	"protractor":       {i: "\U0000f288", c: [3]uint8{229, 61, 58}},   // protractor (Not supported by nerdFont)
	"pug":              {i: "\U0000e686", c: [3]uint8{239, 204, 163}}, // pug
	"puppet":           {i: "\U0000e631", c: [3]uint8{251, 193, 60}},  // puppet
	"purescript":       {i: "\U0000e630", c: [3]uint8{66, 165, 245}},  // purescript
	"python":           {i: "\U000f0320", c: [3]uint8{52, 102, 143}},  // python
	"python-misc":      {i: "\U0000f820", c: [3]uint8{130, 61, 28}},   // python-misc
	"qsharp":           {i: "\U0000f292", c: [3]uint8{251, 193, 60}},  // qsharp (Not supported by nerdFont)
	"r":                {i: "\U000f07d4", c: [3]uint8{25, 118, 210}},  // r
	"raml":             {i: "\U0000e60b", c: [3]uint8{66, 165, 245}},  // raml
	"raspberry-pi":     {i: "\U0000f315", c: [3]uint8{208, 60, 76}},   // raspberry-pi
	"razor":            {i: "\U000f0065", c: [3]uint8{66, 165, 245}},  // razor
	"react":            {i: "\U0000e7ba", c: [3]uint8{35, 188, 212}},  // react
	"react_ts":         {i: "\U0000e7ba", c: [3]uint8{36, 142, 211}},  // react_ts
	"readme":           {i: "\U000f02fc", c: [3]uint8{66, 165, 245}},  // readme
	"redhat":           {i: "\U0000f316", c: [3]uint8{231, 61, 58}},   // redhat
	"roadmap":          {i: "\U000f066e", c: [3]uint8{48, 166, 154}},  // roadmap
	"robot":            {i: "\U000f06a9", c: [3]uint8{249, 89, 63}},   // robot
	"rollup":           {i: "\U000f0bc0", c: [3]uint8{233, 61, 50}},   // rollup
	"routing":          {i: "\U000f0641", c: [3]uint8{67, 160, 71}},   // routing
	"ruby":             {i: "\U0000e739", c: [3]uint8{229, 61, 58}},   // ruby
	"rust":             {i: "\U0000e7a8", c: [3]uint8{250, 111, 66}},  // rust
	"sass":             {i: "\U0000e603", c: [3]uint8{237, 80, 122}},  // sass
	"scheme":           {i: "\U000f0627", c: [3]uint8{244, 68, 62}},   // scheme
	"semantic-release": {i: "\U000f0210", c: [3]uint8{245, 245, 245}}, // semantic-release (Not supported by nerdFont)
	"settings":         {i: "\U0000f013", c: [3]uint8{66, 165, 245}},  // settings
	"shaderlab":        {i: "\U000f06af", c: [3]uint8{25, 118, 210}},  // shaderlab
	"sketch":           {i: "\U000f0b8a", c: [3]uint8{255, 194, 61}},  // sketch
	"slim":             {i: "\U0000f24e", c: [3]uint8{245, 129, 61}},  // slim (Not supported by nerdFont)
	"smarty":           {i: "\U000f0335", c: [3]uint8{255, 207, 60}},  // smarty
	"solidity":         {i: "\U000f07bb", c: [3]uint8{3, 136, 209}},   // solidity
	"sqlite":           {i: "\U0000e7c4", c: [3]uint8{1, 57, 84}},     // sqlite
	"storybook":        {i: "\U000f082e", c: [3]uint8{237, 80, 122}},  // storybook (Not supported by nerdFont)
	"stryker":          {i: "\U0000f05b", c: [3]uint8{239, 83, 80}},   // stryker
	"stylelint":        {i: "\U0000e695", c: [3]uint8{207, 216, 220}}, // stylelint
	"stylus":           {i: "\U0000e600", c: [3]uint8{192, 202, 51}},  // stylus
	"sublime":          {i: "\U0000e696", c: [3]uint8{239, 148, 58}},  // sublime
	"svelte":           {i: "\U0000e697", c: [3]uint8{255, 62, 0}},    // svelte
	"svg":              {i: "\U000f0721", c: [3]uint8{255, 181, 62}},  // svg
	"swc":              {i: "\U000f06d5", c: [3]uint8{198, 52, 54}},   // swc (Not supported by nerdFont)
	"swift":            {i: "\U000f06e5", c: [3]uint8{249, 95, 63}},   // swift
	"table":            {i: "\U000f021b", c: [3]uint8{139, 195, 74}},  // table
	"tailwindcss":      {i: "\U000f13ff", c: [3]uint8{77, 182, 172}},  // tailwindcss
	"tcl":              {i: "\U000f06d3", c: [3]uint8{239, 83, 80}},   // tcl
	"terraform":        {i: "\U000f1062", c: [3]uint8{92, 107, 192}},  // terraform
	"test-js":          {i: "\U000f0096", c: [3]uint8{255, 202, 61}},  // test-js
	"test-jsx":         {i: "\U000f0096", c: [3]uint8{35, 188, 212}},  // test-jsx
	"test-ts":          {i: "\U000f0096", c: [3]uint8{3, 136, 209}},   // test-ts
	"tex":              {i: "\U0000e69b", c: [3]uint8{66, 165, 245}},  // tex
	"todo":             {i: "\U0000f058", c: [3]uint8{124, 179, 66}},  // todo
	"travis":           {i: "\U0000e77e", c: [3]uint8{203, 58, 73}},   // travis
	"tune":             {i: "\U000f066a", c: [3]uint8{251, 193, 60}},  // tune
	"twig":             {i: "\U0000e61c", c: [3]uint8{155, 185, 47}},  // twig
	"typescript":       {i: "\U0000e628", c: [3]uint8{3, 136, 209}},   // typescript
	"typescript-def":   {i: "\U000f06e6", c: [3]uint8{3, 136, 209}},   // typescript-def
	"typst":            {i: "t", c: [3]uint8{34, 157, 172}},           // typst
	"ubuntu":           {i: "\U0000f31c", c: [3]uint8{214, 73, 53}},   // ubuntu
	"url":              {i: "\U000f0337", c: [3]uint8{66, 165, 245}},  // url
	"vagrant":          {i: "\U0000f27d", c: [3]uint8{20, 101, 192}},  // vagrant (Not supported by nerdFont)
	"vala":             {i: "\U0000e69e", c: [3]uint8{149, 117, 205}}, // vala (Not supported by nerdFont)
	"vercel":           {i: "\U0000f47e", c: [3]uint8{207, 216, 220}}, // vercel
	"verilog":          {i: "\U000f061a", c: [3]uint8{250, 111, 66}},  // verilog
	"video":            {i: "\U000f022b", c: [3]uint8{253, 154, 62}},  // video
	"vim":              {i: "\U0000e62b", c: [3]uint8{67, 160, 71}},   // vim
	"virtual":          {i: "\U0000f822", c: [3]uint8{3, 155, 229}},   // virtual
	"visualstudio":     {i: "\U0000e70c", c: [3]uint8{173, 99, 188}},  // visualstudio
	"vscode":           {i: "\U0000e70c", c: [3]uint8{33, 150, 243}},  // vscode (Not supported by nerdFont)
	"vue":              {i: "\U000f0844", c: [3]uint8{65, 184, 131}},  // vue
	"vue-config":       {i: "\U000f0844", c: [3]uint8{58, 121, 110}},  // vue-config
	"webpack":          {i: "\U000f072b", c: [3]uint8{142, 214, 251}}, // webpack
	"word":             {i: "\U000f022c", c: [3]uint8{1, 87, 155}},    // word
	"xaml":             {i: "\U000f05c0", c: [3]uint8{64, 153, 69}},   // xaml
	"xml":              {i: "\U000f05c0", c: [3]uint8{64, 153, 69}},   // xml
	"yaml":             {i: "\U0000e60b", c: [3]uint8{244, 68, 62}},   // yaml
	"yang":             {i: "\U000f0680", c: [3]uint8{66, 165, 245}},  // yang
	"yarn":             {i: "\U000f011b", c: [3]uint8{44, 142, 187}},  // yarn
	"zig":              {i: "\U0000e6a9", c: [3]uint8{249, 169, 60}},  // zig
	"zip":              {i: "\U0000f410", c: [3]uint8{175, 180, 43}},  // zip

	// "abc":              {i:"\u", c:[3]uint8{255, 255, 255}},       // abc
	// "adonis":           {i:"\u", c:[3]uint8{255, 255, 255}},       // adonis
	// "advpl_include":    {i:"\u", c:[3]uint8{255, 255, 255}},       // advpl_include
	// "advpl_prw":        {i:"\u", c:[3]uint8{255, 255, 255}},       // advpl_prw
	// "advpl_ptm":        {i:"\u", c:[3]uint8{255, 255, 255}},       // advpl_ptm
	// "advpl_tlpp":       {i:"\u", c:[3]uint8{255, 255, 255}},       // advpl_tlpp
	// "apollo":           {i:"\u", c:[3]uint8{255, 255, 255}},       // apollo
	// "appveyor":         {i:"\u", c:[3]uint8{255, 255, 255}},       // appveyor
	// "aurelia":          {i:"\u", c:[3]uint8{255, 255, 255}},       // aurelia
	// "autoit":           {i:"\u", c:[3]uint8{255, 255, 255}},       // autoit
	// "ballerina":        {i:"\u", c:[3]uint8{255, 255, 255}},       // ballerina
	// "bazel":            {i:"\u", c:[3]uint8{255, 255, 255}},       // bazel
	// "bithound":         {i:"\u", c:[3]uint8{255, 255, 255}},       // bithound
	// "browserlist":      {i:"\u", c:[3]uint8{255, 255, 255}},       // browserlist
	// "buck":             {i:"\u", c:[3]uint8{255, 255, 255}},       // buck
	// "bucklescript":     {i:"\u", c:[3]uint8{255, 255, 255}},       // bucklescript
	// "buildkite":        {i:"\u", c:[3]uint8{255, 255, 255}},       // buildkite
	// "cabal":            {i:"\u", c:[3]uint8{255, 255, 255}},       // cabal
	// "capacitor":        {i:"\u", c:[3]uint8{255, 255, 255}},       // capacitor
	// "circleci":         {i:"\u", c:[3]uint8{255, 255, 255}},       // circleci
	// "cloudfoundry":     {i:"\u", c:[3]uint8{255, 255, 255}},       // cloudfoundry
	// "coldfusion":       {i:"\u", c:[3]uint8{255, 255, 255}},       // coldfusion
	// "crystal":          {i:"\u", c:[3]uint8{255, 255, 255}},       // crystal
	// "cucumber":         {i:"\u", c:[3]uint8{255, 255, 255}},       // cucumber
	// "cuda":             {i:"\u", c:[3]uint8{255, 255, 255}},       // cuda
	// "dotjs":            {i:"\u", c:[3]uint8{255, 255, 255}},       // dotjs
	// "drawio":           {i:"\u", c:[3]uint8{255, 255, 255}},       // drawio
	// "drone":            {i:"\u", c:[3]uint8{255, 255, 255}},       // drone
	// "editorconfig":     {i:"\u", c:[3]uint8{255, 255, 255}},       // editorconfig
	// "flow":             {i:"\u", c:[3]uint8{255, 255, 255}},       // flow
	// "forth":            {i:"\u", c:[3]uint8{255, 255, 255}},       // forth
	// "foxpro":           {i:"\u", c:[3]uint8{255, 255, 255}},       // foxpro
	// "fusebox":          {i:"\u", c:[3]uint8{255, 255, 255}},       // fusebox
	// "gitpod":           {i:"\u", c:[3]uint8{255, 255, 255}},       // gitpod
	// "graphcool":        {i:"\u", c:[3]uint8{255, 255, 255}},       // graphcool
	// "hack":             {i:"\u", c:[3]uint8{255, 255, 255}},       // hack
	// "haml":             {i:"\u", c:[3]uint8{255, 255, 255}},       // haml
	// "hcl":              {i:"\u", c:[3]uint8{255, 255, 255}},       // hcl
	// "imba":             {i:"\u", c:[3]uint8{255, 255, 255}},       // imba
	// "istanbul":         {i:"\u", c:[3]uint8{255, 255, 255}},       // istanbul
	// "jupyter":          {i:"\u", c:[3]uint8{255, 255, 255}},       // jupyter
	// "kivy":             {i:"\u", c:[3]uint8{255, 255, 255}},       // kivy
	// "kl":               {i:"\u", c:[3]uint8{255, 255, 255}},       // kl
	// "lisp":             {i:"\u", c:[3]uint8{255, 255, 255}},       // lisp
	// "livescript":       {i:"\u", c:[3]uint8{255, 255, 255}},       // livescript
	// "mathematica":      {i:"\u", c:[3]uint8{255, 255, 255}},       // mathematica
	// "meson":            {i:"\u", c:[3]uint8{255, 255, 255}},       // meson
	// "mint":             {i:"\u", c:[3]uint8{255, 255, 255}},       // mint
	// "nest":             {i:"\u", c:[3]uint8{255, 255, 255}},       // nest
	// "netlify":          {i:"\u", c:[3]uint8{255, 255, 255}},       // netlify
	// "nodemon":          {i:"\u", c:[3]uint8{255, 255, 255}},       // nodemon
	// "nrwl":             {i:"\u", c:[3]uint8{255, 255, 255}},       // nrwl
	// "nunjucks":         {i:"\u", c:[3]uint8{255, 255, 255}},       // nunjucks
	// "percy":            {i:"\u", c:[3]uint8{255, 255, 255}},       // percy
	// "processing":       {i:"\u", c:[3]uint8{255, 255, 255}},       // processing
	// "racket":           {i:"\u", c:[3]uint8{255, 255, 255}},       // racket
	// "reason":           {i:"\u", c:[3]uint8{255, 255, 255}},       // reason
	// "red":              {i:"\u", c:[3]uint8{255, 255, 255}},       // red
	// "restql":           {i:"\u", c:[3]uint8{255, 255, 255}},       // restql
	// "riot":             {i:"\u", c:[3]uint8{255, 255, 255}},       // riot
	// "san":              {i:"\u", c:[3]uint8{255, 255, 255}},       // san
	// "sas":              {i:"\u", c:[3]uint8{255, 255, 255}},       // sas
	// "sbt":              {i:"\u", c:[3]uint8{255, 255, 255}},       // sbt
	// "sequelize":        {i:"\u", c:[3]uint8{255, 255, 255}},       // sequelize
	// "slug":             {i:"\u", c:[3]uint8{255, 255, 255}},       // slug
	// "sml":              {i:"\u", c:[3]uint8{255, 255, 255}},       // sml
	// "snyk":             {i:"\u", c:[3]uint8{255, 255, 255}},       // snyk
	// "stencil":          {i:"\u", c:[3]uint8{255, 255, 255}},       // stencil
	// "tilt":             {i:"\u", c:[3]uint8{255, 255, 255}},       // tilt
	// "uml":              {i:"\u", c:[3]uint8{255, 255, 255}},       // uml
	// "velocity":         {i:"\u", c:[3]uint8{255, 255, 255}},       // velocity
	// "vfl":              {i:"\u", c:[3]uint8{255, 255, 255}},       // vfl
	// "wakatime":         {i:"\u", c:[3]uint8{255, 255, 255}},       // wakatime
	// "wallaby":          {i:"\u", c:[3]uint8{255, 255, 255}},       // wallaby
	// "watchman":         {i:"\u", c:[3]uint8{255, 255, 255}},       // watchman
	// "webassembly":      {i:"\u", c:[3]uint8{255, 255, 255}},       // webassembly
	// "webhint":          {i:"\u", c:[3]uint8{255, 255, 255}},       // webhint
	// "wepy":             {i:"\u", c:[3]uint8{255, 255, 255}},       // wepy
	// "wolframlanguage":  {i:"\u", c:[3]uint8{255, 255, 255}},       // wolframlanguage

	"dir-config":      {i: "\U0000e5fc", c: [3]uint8{32, 173, 194}},  // dir-config
	"dir-controller":  {i: "\U0000e5fc", c: [3]uint8{255, 194, 61}},  // dir-controller
	"dir-download":    {i: "\U000f024d", c: [3]uint8{76, 175, 80}},   // dir-download
	"dir-environment": {i: "\U000f024f", c: [3]uint8{102, 187, 106}}, // dir-environment
	"dir-git":         {i: "\U0000e5fb", c: [3]uint8{250, 111, 66}},  // dir-git
	"dir-github":      {i: "\U0000e5fd", c: [3]uint8{84, 110, 122}},  // dir-github
	"dir-images":      {i: "\U000f024f", c: [3]uint8{43, 150, 137}},  // dir-images
	"dir-import":      {i: "\U000f0257", c: [3]uint8{175, 180, 43}},  // dir-import
	"dir-include":     {i: "\U000f0257", c: [3]uint8{3, 155, 229}},   // dir-include
	"dir-npm":         {i: "\U0000e5fa", c: [3]uint8{203, 56, 55}},   // dir-npm
	"dir-secure":      {i: "\U000f0250", c: [3]uint8{249, 169, 60}},  // dir-secure
	"dir-upload":      {i: "\U000f0259", c: [3]uint8{250, 111, 66}},  // dir-upload
}

// default icons in case nothing can be found
var IconDef = map[string]*IconInfo{
	"dir":        {i: "\U000f024b", c: [3]uint8{224, 177, 77}},
	"diropen":    {i: "\U000f0770", c: [3]uint8{224, 177, 77}},
	"exe":        {i: "\U000f0214", c: [3]uint8{76, 175, 80}},
	"file":       {i: "\U000f0224", c: [3]uint8{65, 129, 190}},
	"hiddendir":  {i: "\U000f0256", c: [3]uint8{224, 177, 77}},
	"hiddenfile": {i: "\U000f0613", c: [3]uint8{65, 129, 190}},
}
