// assets contains all the Icon glyphs info
package icons

import "fmt"

// Icon_Info (icon information)
type IconInfo struct {
	glyph        string
	color        [3]uint8 // represents the color in rgb (default 0,0,0 is black)
	isExecutable bool     // whether or not the file is executable [true = is executable]
}

func (i *IconInfo) GetGlyph() string {
	return i.glyph
}

func (i *IconInfo) GetColor() string {
	if i.isExecutable {
		return "\033[38;2;76;175;080m"
	}

	return fmt.Sprintf("\033[38;2;%03d;%03d;%03dm", i.color[0], i.color[1], i.color[2])
}

func (i *IconInfo) MakeExe() {
	i.isExecutable = true
}

var IconSet = map[string]*IconInfo{
	"3d":               {glyph: "\U0000e79b", color: [3]uint8{40, 182, 246}},  // 3d
	"actionscript":     {glyph: "\U0000fb25", color: [3]uint8{244, 68, 62}},   // actionscript (Not supported by nerdFont)
	"alpine":           {glyph: "\U0000f300", color: [3]uint8{14, 87, 123}},   // alpine
	"android":          {glyph: "\U000f0032", color: [3]uint8{139, 195, 74}},  // android
	"apiblueprint":     {glyph: "\U0000f031", color: [3]uint8{66, 165, 245}},  // apiblueprint (Not supported by nerdFont)
	"applescript":      {glyph: "\U0000f302", color: [3]uint8{120, 144, 156}}, // applescript
	"arch":             {glyph: "\U0000f303", color: [3]uint8{33, 142, 202}},  // arch
	"arduino":          {glyph: "\U0000f34b", color: [3]uint8{35, 151, 156}},  // arduino
	"asciidoc":         {glyph: "\U000f0219", color: [3]uint8{244, 68, 62}},   // asciidoc
	"assembly":         {glyph: "\U0000f471", color: [3]uint8{250, 109, 63}},  // assembly
	"audio":            {glyph: "\U000f0388", color: [3]uint8{239, 83, 80}},   // audio
	"authors":          {glyph: "\U0000f0c0", color: [3]uint8{244, 68, 62}},   // authors
	"autohotkey":       {glyph: "\U000f0313", color: [3]uint8{76, 175, 80}},   // autohotkey (Not supported by nerdFont)
	"azure":            {glyph: "\U0000fd03", color: [3]uint8{31, 136, 229}},  // azure
	"azure-pipelines":  {glyph: "\U0000f427", color: [3]uint8{20, 101, 192}},  // azure-pipelines
	"babel":            {glyph: "\U000f0a25", color: [3]uint8{253, 217, 59}},  // babel
	"bitbucket":        {glyph: "\U0000f171", color: [3]uint8{31, 136, 229}},  // bitbucket
	"blink":            {glyph: "\U000f022b", color: [3]uint8{249, 169, 60}},  // blink (The Foundry Nuke) (Not supported by nerdFont)
	"bower":            {glyph: "\U0000e61a", color: [3]uint8{239, 88, 60}},   // bower
	"bun":              {glyph: "\U000f0685", color: [3]uint8{249, 241, 225}}, // bun
	"c":                {glyph: "\U000f0671", color: [3]uint8{2, 119, 189}},   // c
	"cake":             {glyph: "\U000f00eb", color: [3]uint8{250, 111, 66}},  // cake
	"centOS":           {glyph: "\U0000f304", color: [3]uint8{157, 83, 135}},  // centOS
	"certificate":      {glyph: "\U000f0124", color: [3]uint8{249, 89, 63}},   // certificate
	"changelog":        {glyph: "\U000f099b", color: [3]uint8{139, 195, 74}},  // changelog
	"clojure":          {glyph: "\U0000e76a", color: [3]uint8{100, 221, 23}},  // clojure
	"cmake":            {glyph: "\U0000f425", color: [3]uint8{178, 178, 179}}, // cmake (Not supported by nerdFont)
	"coconut":          {glyph: "\U000f00d3", color: [3]uint8{141, 110, 99}},  // coconut
	"code-climate":     {glyph: "\U000f0509", color: [3]uint8{238, 238, 238}}, // code-climate
	"codecov":          {glyph: "\U0000e37c", color: [3]uint8{237, 80, 122}},  // codecov (Not supported by nerdFont)
	"codeowners":       {glyph: "\U000f0008", color: [3]uint8{175, 180, 43}},  // codeowners
	"coffee":           {glyph: "\U000f0176", color: [3]uint8{66, 165, 245}},  // coffee
	"command":          {glyph: "\U000f0633", color: [3]uint8{175, 188, 194}}, // command
	"commitlint":       {glyph: "\U000f0718", color: [3]uint8{43, 150, 137}},  // commitlint
	"conduct":          {glyph: "\U000f014e", color: [3]uint8{205, 220, 57}},  // conduct
	"console":          {glyph: "\U000f018d", color: [3]uint8{250, 111, 66}},  // console
	"contributing":     {glyph: "\U000f014d", color: [3]uint8{255, 202, 61}},  // contributing
	"cpp":              {glyph: "\U000f0672", color: [3]uint8{2, 119, 189}},   // cpp
	"credits":          {glyph: "\U000f0260", color: [3]uint8{156, 204, 101}}, // credits
	"csharp":           {glyph: "\U000f031b", color: [3]uint8{2, 119, 189}},   // csharp
	"css":              {glyph: "\U000f031c", color: [3]uint8{66, 165, 245}},  // css
	"css-map":          {glyph: "\U0000e749", color: [3]uint8{66, 165, 245}},  // css-map
	"d":                {glyph: "\U0000e7af", color: [3]uint8{244, 68, 62}},   // d
	"dart":             {glyph: "\U0000e798", color: [3]uint8{87, 182, 240}},  // dart
	"database":         {glyph: "\U0000e706", color: [3]uint8{255, 202, 61}},  // database
	"debian":           {glyph: "\U0000f306", color: [3]uint8{211, 61, 76}},   // debian
	"denizenscript":    {glyph: "\U000f0af1", color: [3]uint8{255, 213, 79}},  // denizenscript (Not supported by nerdFont)
	"dhall":            {glyph: "\U0000f448", color: [3]uint8{120, 144, 156}}, // dhall
	"disc":             {glyph: "\U0000e271", color: [3]uint8{176, 190, 197}}, // disc
	"django":           {glyph: "\U0000e71d", color: [3]uint8{67, 160, 71}},   // django
	"docker":           {glyph: "\U0000f308", color: [3]uint8{1, 135, 201}},   // docker
	"document":         {glyph: "\U000f0219", color: [3]uint8{66, 165, 245}},  // document
	"dune":             {glyph: "\U000f0509", color: [3]uint8{244, 127, 61}},  // dune
	"edge":             {glyph: "\U000f0065", color: [3]uint8{239, 111, 60}},  // edge
	"editorconfig":     {glyph: "\U0000e652", color: [3]uint8{255, 242, 242}}, // editorconfig
	"ejs":              {glyph: "\U0000e618", color: [3]uint8{255, 202, 61}},  // ejs
	"elixir":           {glyph: "\U0000e62d", color: [3]uint8{149, 117, 205}}, // elixir
	"elm":              {glyph: "\U0000e62c", color: [3]uint8{96, 181, 204}},  // elm
	"email":            {glyph: "\U000f01ee", color: [3]uint8{66, 165, 245}},  // email
	"erlang":           {glyph: "\U0000e7b1", color: [3]uint8{244, 68, 62}},   // erlang
	"eslint":           {glyph: "\U000f0c7a", color: [3]uint8{121, 134, 203}}, // eslint
	"exe":              {glyph: "\U0000f2d0", color: [3]uint8{229, 77, 58}},   // exe
	"fastlane":         {glyph: "\U000f0700", color: [3]uint8{149, 119, 232}}, // fastlane (Not supported by nerdFont)
	"favicon":          {glyph: "\U0000e623", color: [3]uint8{255, 213, 79}},  // favicon
	"fedora":           {glyph: "\U0000f30a", color: [3]uint8{52, 103, 172}},  // fedora
	"firebase":         {glyph: "\U000f0967", color: [3]uint8{251, 193, 60}},  // firebase
	"flash":            {glyph: "\U000f0241", color: [3]uint8{198, 52, 54}},   // flash (Not supported by nerdFont)
	"font":             {glyph: "\U0000f031", color: [3]uint8{244, 68, 62}},   // font
	"fortran":          {glyph: "\U000f121a", color: [3]uint8{250, 111, 66}},  // fortran
	"freebsd":          {glyph: "\U0000f30c", color: [3]uint8{175, 44, 42}},   // freebsd
	"fsharp":           {glyph: "\U0000e7a7", color: [3]uint8{55, 139, 186}},  // fsharp
	"gatsby":           {glyph: "\U000f0e43", color: [3]uint8{70, 0, 130}},    // gatsby
	"gcp":              {glyph: "\U000f0163", color: [3]uint8{70, 136, 250}},  // gcp (Not supported by nerdFont)
	"gemfile":          {glyph: "\U0000e21e", color: [3]uint8{229, 61, 58}},   // gemfile
	"gentoo":           {glyph: "\U0000f30d", color: [3]uint8{148, 141, 211}}, // gentoo
	"gnu":              {glyph: "\U0000e779", color: [3]uint8{229, 61, 58}},   // GNU
	"git":              {glyph: "\U0000e702", color: [3]uint8{229, 77, 58}},   // git
	"gitlab":           {glyph: "\U0000f296", color: [3]uint8{226, 69, 57}},   // gitlab
	"go":               {glyph: "\U000f07d3", color: [3]uint8{32, 173, 194}},  // go
	"go-mod":           {glyph: "\U000f07d3", color: [3]uint8{237, 80, 122}},  // go-mod
	"go-test":          {glyph: "\U000f07d3", color: [3]uint8{255, 213, 79}},  // go-test
	"godot":            {glyph: "\U0000e65f", color: [3]uint8{79, 195, 247}},  // godot
	"godot-assets":     {glyph: "\U0000e65f", color: [3]uint8{129, 199, 132}}, // godot-assets
	"gradle":           {glyph: "\U0000e660", color: [3]uint8{29, 151, 167}},  // gradle
	"graphql":          {glyph: "\U000f0877", color: [3]uint8{237, 80, 122}},  // graphql
	"groovy":           {glyph: "\U0000f2a6", color: [3]uint8{41, 198, 218}},  // groovy
	"grunt":            {glyph: "\U0000e611", color: [3]uint8{251, 170, 61}},  // grunt
	"gulp":             {glyph: "\U0000e763", color: [3]uint8{229, 61, 58}},   // gulp
	"h":                {glyph: "\U000f0c00", color: [3]uint8{2, 119, 189}},   // h (Not supported by nerdFont)
	"handlebars":       {glyph: "\U0000e60f", color: [3]uint8{250, 111, 66}},  // handlebars
	"haskell":          {glyph: "\U0000e61f", color: [3]uint8{254, 168, 62}},  // haskell
	"haxe":             {glyph: "\U0000e666", color: [3]uint8{246, 137, 61}},  // haxe
	"helm":             {glyph: "\U000f0833", color: [3]uint8{32, 173, 194}},  // helm (Not supported by nerdFont)
	"heroku":           {glyph: "\U0000e607", color: [3]uint8{105, 99, 185}},  // heroku
	"hpp":              {glyph: "\U000f0b0f", color: [3]uint8{2, 119, 189}},   // hpp (Not supported by nerdFont)
	"html":             {glyph: "\U0000f13b", color: [3]uint8{228, 79, 57}},   // html
	"http":             {glyph: "\U0000f484", color: [3]uint8{66, 165, 245}},  // http
	"husky":            {glyph: "\U000f03e9", color: [3]uint8{229, 229, 229}}, // husky
	"i18n":             {glyph: "\U000f05ca", color: [3]uint8{121, 134, 203}}, // i18n (Not supported by nerdFont)
	"image":            {glyph: "\U000f021f", color: [3]uint8{48, 166, 154}},  // image
	"ionic":            {glyph: "\U0000e7a9", color: [3]uint8{79, 143, 247}},  // ionic
	"java":             {glyph: "\U000f0176", color: [3]uint8{244, 68, 62}},   // java
	"javascript":       {glyph: "\U0000e74e", color: [3]uint8{255, 202, 61}},  // javascript
	"javascript-map":   {glyph: "\U0000e781", color: [3]uint8{255, 202, 61}},  // javascript-map
	"jenkins":          {glyph: "\U0000e767", color: [3]uint8{240, 214, 183}}, // jenkins
	"jest":             {glyph: "\U000f0af7", color: [3]uint8{244, 85, 62}},   // jest (Not supported by nerdFont)
	"jinja":            {glyph: "\U0000e66f", color: [3]uint8{174, 44, 42}},   // jinja
	"json":             {glyph: "\U0000e60b", color: [3]uint8{251, 193, 60}},  // json
	"julia":            {glyph: "\U0000e624", color: [3]uint8{134, 82, 159}},  // julia
	"karma":            {glyph: "\U0000e622", color: [3]uint8{60, 190, 174}},  // karma
	"key":              {glyph: "\U000f0306", color: [3]uint8{48, 166, 154}},  // key
	"kotlin":           {glyph: "\U0000e634", color: [3]uint8{139, 195, 74}},  // kotlin
	"laravel":          {glyph: "\U0000e73f", color: [3]uint8{248, 80, 81}},   // laravel
	"less":             {glyph: "\U0000e60b", color: [3]uint8{2, 119, 189}},   // less
	"lib":              {glyph: "\U0000f831", color: [3]uint8{139, 195, 74}},  // lib
	"linux":            {glyph: "\U0000e712", color: [3]uint8{238, 207, 55}},  // linux
	"liquid":           {glyph: "\U0000e275", color: [3]uint8{40, 182, 246}},  // liquid
	"lock":             {glyph: "\U000f033e", color: [3]uint8{255, 213, 79}},  // lock
	"log":              {glyph: "\U0000f719", color: [3]uint8{175, 180, 43}},  // log
	"lua":              {glyph: "\U0000e620", color: [3]uint8{66, 165, 245}},  // lua
	"makefile":         {glyph: "\U000f0229", color: [3]uint8{239, 83, 80}},   // makefile
	"manjaro":          {glyph: "\U0000f312", color: [3]uint8{73, 185, 90}},   // manjaro
	"markdown":         {glyph: "\U0000eb1d", color: [3]uint8{66, 165, 245}},  // markdown
	"markojs":          {glyph: "\U0000f13b", color: [3]uint8{2, 119, 189}},   // markojs (Not supported by nerdFont)
	"mdx":              {glyph: "\U0000f853", color: [3]uint8{255, 202, 61}},  // mdx
	"merlin":           {glyph: "\U0000f136", color: [3]uint8{66, 165, 245}},  // merlin (Not supported by nerdFont)
	"midi":             {glyph: "\U000f08f2", color: [3]uint8{66, 165, 245}},  // midi
	"mint":             {glyph: "\U0000f30f", color: [3]uint8{125, 190, 58}},  // mint
	"mjml":             {glyph: "\U0000e714", color: [3]uint8{249, 89, 63}},   // mjml (Not supported by nerdFont)
	"mocha":            {glyph: "\U000f01aa", color: [3]uint8{161, 136, 127}}, // mocha (Not supported by nerdFont)
	"modernizr":        {glyph: "\U0000e720", color: [3]uint8{234, 72, 99}},   // modernizr
	"moonscript":       {glyph: "\U0000f186", color: [3]uint8{251, 193, 60}},  // moonscript
	"mxml":             {glyph: "\U0000f72d", color: [3]uint8{254, 168, 62}},  // mxml
	"mysql":            {glyph: "\U0000e704", color: [3]uint8{1, 94, 134}},    // mysql
	"nim":              {glyph: "\U000f01a5", color: [3]uint8{255, 202, 61}},  // nim
	"nix":              {glyph: "\U0000f313", color: [3]uint8{80, 117, 193}},  // nix
	"nodejs":           {glyph: "\U000f0399", color: [3]uint8{139, 195, 74}},  // nodejs
	"npm":              {glyph: "\U0000e71e", color: [3]uint8{203, 56, 55}},   // npm
	"nuget":            {glyph: "\U0000e77f", color: [3]uint8{3, 136, 209}},   // nuget (Not supported by nerdFont)
	"nuxt":             {glyph: "\U000f1106", color: [3]uint8{65, 184, 131}},  // nuxt
	"ocaml":            {glyph: "\U0000e67a", color: [3]uint8{253, 154, 62}},  // ocaml
	"opam":             {glyph: "\U0000f1ce", color: [3]uint8{255, 213, 79}},  // opam (Not supported by nerdFont)
	"opensuse":         {glyph: "\U0000f314", color: [3]uint8{111, 180, 36}},  // opensuse
	"pascal":           {glyph: "\U000f03db", color: [3]uint8{3, 136, 209}},   // pascal (Not supported by nerdFont)
	"pawn":             {glyph: "\U0000e261", color: [3]uint8{239, 111, 60}},  // pawn
	"pdf":              {glyph: "\U0000e67d", color: [3]uint8{244, 68, 62}},   // pdf
	"perl":             {glyph: "\U0000e769", color: [3]uint8{149, 117, 205}}, // perl
	"php":              {glyph: "\U0000e608", color: [3]uint8{65, 129, 190}},  // php
	"postcss":          {glyph: "\U000f031c", color: [3]uint8{244, 68, 62}},   // postcss (Not supported by nerdFont)
	"postgresql":       {glyph: "\U0000e76e", color: [3]uint8{49, 99, 140}},   // postgresql
	"powerpoint":       {glyph: "\U000f0227", color: [3]uint8{209, 71, 51}},   // powerpoint
	"powershell":       {glyph: "\U000f0a0a", color: [3]uint8{5, 169, 244}},   // powershell
	"prettier":         {glyph: "\U0000e6b4", color: [3]uint8{86, 179, 180}},  // prettier
	"prisma":           {glyph: "\U0000e684", color: [3]uint8{229, 56, 150}},  // prisma
	"prolog":           {glyph: "\U0000e7a1", color: [3]uint8{239, 83, 80}},   // prolog
	"protractor":       {glyph: "\U0000f288", color: [3]uint8{229, 61, 58}},   // protractor (Not supported by nerdFont)
	"pug":              {glyph: "\U0000e686", color: [3]uint8{239, 204, 163}}, // pug
	"puppet":           {glyph: "\U0000e631", color: [3]uint8{251, 193, 60}},  // puppet
	"purescript":       {glyph: "\U0000e630", color: [3]uint8{66, 165, 245}},  // purescript
	"python":           {glyph: "\U000f0320", color: [3]uint8{52, 102, 143}},  // python
	"python-misc":      {glyph: "\U0000f820", color: [3]uint8{130, 61, 28}},   // python-misc
	"qsharp":           {glyph: "\U0000f292", color: [3]uint8{251, 193, 60}},  // qsharp (Not supported by nerdFont)
	"r":                {glyph: "\U000f07d4", color: [3]uint8{25, 118, 210}},  // r
	"raml":             {glyph: "\U0000e60b", color: [3]uint8{66, 165, 245}},  // raml
	"raspberry-pi":     {glyph: "\U0000f315", color: [3]uint8{208, 60, 76}},   // raspberry-pi
	"razor":            {glyph: "\U000f0065", color: [3]uint8{66, 165, 245}},  // razor
	"react":            {glyph: "\U0000e7ba", color: [3]uint8{35, 188, 212}},  // react
	"react_ts":         {glyph: "\U0000e7ba", color: [3]uint8{36, 142, 211}},  // react_ts
	"readme":           {glyph: "\U000f02fc", color: [3]uint8{66, 165, 245}},  // readme
	"redhat":           {glyph: "\U0000f316", color: [3]uint8{231, 61, 58}},   // redhat
	"roadmap":          {glyph: "\U000f066e", color: [3]uint8{48, 166, 154}},  // roadmap
	"robot":            {glyph: "\U000f06a9", color: [3]uint8{249, 89, 63}},   // robot
	"rollup":           {glyph: "\U000f0bc0", color: [3]uint8{233, 61, 50}},   // rollup
	"routing":          {glyph: "\U000f0641", color: [3]uint8{67, 160, 71}},   // routing
	"ruby":             {glyph: "\U0000e739", color: [3]uint8{229, 61, 58}},   // ruby
	"rust":             {glyph: "\U0000e7a8", color: [3]uint8{250, 111, 66}},  // rust
	"sass":             {glyph: "\U0000e603", color: [3]uint8{237, 80, 122}},  // sass
	"scheme":           {glyph: "\U000f0627", color: [3]uint8{244, 68, 62}},   // scheme
	"semantic-release": {glyph: "\U000f0210", color: [3]uint8{245, 245, 245}}, // semantic-release (Not supported by nerdFont)
	"settings":         {glyph: "\U0000f013", color: [3]uint8{66, 165, 245}},  // settings
	"shaderlab":        {glyph: "\U000f06af", color: [3]uint8{25, 118, 210}},  // shaderlab
	"sketch":           {glyph: "\U000f0b8a", color: [3]uint8{255, 194, 61}},  // sketch
	"slim":             {glyph: "\U0000f24e", color: [3]uint8{245, 129, 61}},  // slim (Not supported by nerdFont)
	"smarty":           {glyph: "\U000f0335", color: [3]uint8{255, 207, 60}},  // smarty
	"solidity":         {glyph: "\U000f07bb", color: [3]uint8{3, 136, 209}},   // solidity
	"sqlite":           {glyph: "\U0000e7c4", color: [3]uint8{1, 57, 84}},     // sqlite
	"storybook":        {glyph: "\U000f082e", color: [3]uint8{237, 80, 122}},  // storybook (Not supported by nerdFont)
	"stryker":          {glyph: "\U0000f05b", color: [3]uint8{239, 83, 80}},   // stryker
	"stylelint":        {glyph: "\U0000e695", color: [3]uint8{207, 216, 220}}, // stylelint
	"stylus":           {glyph: "\U0000e600", color: [3]uint8{192, 202, 51}},  // stylus
	"sublime":          {glyph: "\U0000e696", color: [3]uint8{239, 148, 58}},  // sublime
	"svelte":           {glyph: "\U0000e697", color: [3]uint8{255, 62, 0}},    // svelte
	"svg":              {glyph: "\U000f0721", color: [3]uint8{255, 181, 62}},  // svg
	"swc":              {glyph: "\U000f06d5", color: [3]uint8{198, 52, 54}},   // swc (Not supported by nerdFont)
	"swift":            {glyph: "\U000f06e5", color: [3]uint8{249, 95, 63}},   // swift
	"table":            {glyph: "\U000f021b", color: [3]uint8{139, 195, 74}},  // table
	"tailwindcss":      {glyph: "\U000f13ff", color: [3]uint8{77, 182, 172}},  // tailwindcss
	"tcl":              {glyph: "\U000f06d3", color: [3]uint8{239, 83, 80}},   // tcl
	"terraform":        {glyph: "\U000f1062", color: [3]uint8{92, 107, 192}},  // terraform
	"test-js":          {glyph: "\U000f0096", color: [3]uint8{255, 202, 61}},  // test-js
	"test-jsx":         {glyph: "\U000f0096", color: [3]uint8{35, 188, 212}},  // test-jsx
	"test-ts":          {glyph: "\U000f0096", color: [3]uint8{3, 136, 209}},   // test-ts
	"tex":              {glyph: "\U0000e69b", color: [3]uint8{66, 165, 245}},  // tex
	"todo":             {glyph: "\U0000f058", color: [3]uint8{124, 179, 66}},  // todo
	"travis":           {glyph: "\U0000e77e", color: [3]uint8{203, 58, 73}},   // travis
	"tune":             {glyph: "\U000f066a", color: [3]uint8{251, 193, 60}},  // tune
	"twig":             {glyph: "\U0000e61c", color: [3]uint8{155, 185, 47}},  // twig
	"typescript":       {glyph: "\U0000e628", color: [3]uint8{3, 136, 209}},   // typescript
	"typescript-def":   {glyph: "\U000f06e6", color: [3]uint8{3, 136, 209}},   // typescript-def
	"typst":            {glyph: "t", color: [3]uint8{34, 157, 172}},           // typst
	"ubuntu":           {glyph: "\U0000f31c", color: [3]uint8{214, 73, 53}},   // ubuntu
	"url":              {glyph: "\U000f0337", color: [3]uint8{66, 165, 245}},  // url
	"vagrant":          {glyph: "\U0000f27d", color: [3]uint8{20, 101, 192}},  // vagrant (Not supported by nerdFont)
	"vala":             {glyph: "\U0000e69e", color: [3]uint8{149, 117, 205}}, // vala (Not supported by nerdFont)
	"vercel":           {glyph: "\U0000f47e", color: [3]uint8{207, 216, 220}}, // vercel
	"verilog":          {glyph: "\U000f061a", color: [3]uint8{250, 111, 66}},  // verilog
	"video":            {glyph: "\U000f022b", color: [3]uint8{253, 154, 62}},  // video
	"vim":              {glyph: "\U0000e62b", color: [3]uint8{67, 160, 71}},   // vim
	"virtual":          {glyph: "\U0000f822", color: [3]uint8{3, 155, 229}},   // virtual
	"visualstudio":     {glyph: "\U0000e70c", color: [3]uint8{173, 99, 188}},  // visualstudio
	"vscode":           {glyph: "\U0000e70c", color: [3]uint8{33, 150, 243}},  // vscode (Not supported by nerdFont)
	"vue":              {glyph: "\U000f0844", color: [3]uint8{65, 184, 131}},  // vue
	"vue-config":       {glyph: "\U000f0844", color: [3]uint8{58, 121, 110}},  // vue-config
	"webpack":          {glyph: "\U000f072b", color: [3]uint8{142, 214, 251}}, // webpack
	"word":             {glyph: "\U000f022c", color: [3]uint8{1, 87, 155}},    // word
	"xaml":             {glyph: "\U000f05c0", color: [3]uint8{64, 153, 69}},   // xaml
	"xml":              {glyph: "\U000f05c0", color: [3]uint8{64, 153, 69}},   // xml
	"yaml":             {glyph: "\U0000e60b", color: [3]uint8{244, 68, 62}},   // yaml
	"yang":             {glyph: "\U000f0680", color: [3]uint8{66, 165, 245}},  // yang
	"yarn":             {glyph: "\U000f011b", color: [3]uint8{44, 142, 187}},  // yarn
	"zig":              {glyph: "\U0000e6a9", color: [3]uint8{249, 169, 60}},  // zig
	"zip":              {glyph: "\U0000f410", color: [3]uint8{175, 180, 43}},  // zip

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

	"dir-config":      {glyph: "\U0000e5fc", color: [3]uint8{32, 173, 194}},  // dir-config
	"dir-controller":  {glyph: "\U0000e5fc", color: [3]uint8{255, 194, 61}},  // dir-controller
	"dir-download":    {glyph: "\U000f024d", color: [3]uint8{76, 175, 80}},   // dir-download
	"dir-environment": {glyph: "\U000f024f", color: [3]uint8{102, 187, 106}}, // dir-environment
	"dir-git":         {glyph: "\U0000e5fb", color: [3]uint8{250, 111, 66}},  // dir-git
	"dir-github":      {glyph: "\U0000e5fd", color: [3]uint8{84, 110, 122}},  // dir-github
	"dir-images":      {glyph: "\U000f024f", color: [3]uint8{43, 150, 137}},  // dir-images
	"dir-import":      {glyph: "\U000f0257", color: [3]uint8{175, 180, 43}},  // dir-import
	"dir-include":     {glyph: "\U000f0257", color: [3]uint8{3, 155, 229}},   // dir-include
	"dir-npm":         {glyph: "\U0000e5fa", color: [3]uint8{203, 56, 55}},   // dir-npm
	"dir-secure":      {glyph: "\U000f0250", color: [3]uint8{249, 169, 60}},  // dir-secure
	"dir-upload":      {glyph: "\U000f0259", color: [3]uint8{250, 111, 66}},  // dir-upload
}

// default icons in case nothing can be found
var IconDef = map[string]*IconInfo{
	"dir":        {glyph: "\U000f024b", color: [3]uint8{224, 177, 77}},
	"diropen":    {glyph: "\U000f0770", color: [3]uint8{224, 177, 77}},
	"exe":        {glyph: "\U000f0214", color: [3]uint8{76, 175, 80}},
	"file":       {glyph: "\U000f0224", color: [3]uint8{65, 129, 190}},
	"hiddendir":  {glyph: "\U000f0256", color: [3]uint8{224, 177, 77}},
	"hiddenfile": {glyph: "\U000f0613", color: [3]uint8{65, 129, 190}},
}
