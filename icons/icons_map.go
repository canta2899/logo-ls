// assets contains all the Icon glyphs info
package icons

import "fmt"

// Icon_Info (icon information)
type IconInfo struct {
	Glyph        string
	Color        [3]uint8 // represents the color in rgb (default 0,0,0 is black)
	IsExecutable bool     // whether or not the file is executable [true = is executable]
}

func (i *IconInfo) GetGlyph() string {
	if i == nil {
		return ""
	}

	return i.Glyph
}

func (i *IconInfo) GetColor() string {
	if i == nil {
		return ""
	}

	if i.IsExecutable {
		return "\033[38;2;76;175;080m"
	}

	return fmt.Sprintf("\033[38;2;%03d;%03d;%03dm", i.Color[0], i.Color[1], i.Color[2])
}

func (i *IconInfo) MakeExe() {
	i.IsExecutable = true
}

var IconSet = map[string]*IconInfo{
	"3d":               {Glyph: "\U0000e79b", Color: [3]uint8{40, 182, 246}},  // 3d
	"actionscript":     {Glyph: "\U0000fb25", Color: [3]uint8{244, 68, 62}},   // actionscript (Not supported by nerdFont)
	"alpine":           {Glyph: "\U0000f300", Color: [3]uint8{14, 87, 123}},   // alpine
	"android":          {Glyph: "\U000f0032", Color: [3]uint8{139, 195, 74}},  // android
	"apiblueprint":     {Glyph: "\U0000f031", Color: [3]uint8{66, 165, 245}},  // apiblueprint (Not supported by nerdFont)
	"applescript":      {Glyph: "\U0000f302", Color: [3]uint8{120, 144, 156}}, // applescript
	"arch":             {Glyph: "\U0000f303", Color: [3]uint8{33, 142, 202}},  // arch
	"arduino":          {Glyph: "\U0000f34b", Color: [3]uint8{35, 151, 156}},  // arduino
	"asciidoc":         {Glyph: "\U000f0219", Color: [3]uint8{244, 68, 62}},   // asciidoc
	"assembly":         {Glyph: "\U0000f471", Color: [3]uint8{250, 109, 63}},  // assembly
	"audio":            {Glyph: "\U000f0388", Color: [3]uint8{239, 83, 80}},   // audio
	"authors":          {Glyph: "\U0000f0c0", Color: [3]uint8{244, 68, 62}},   // authors
	"autohotkey":       {Glyph: "\U000f0313", Color: [3]uint8{76, 175, 80}},   // autohotkey (Not supported by nerdFont)
	"azure":            {Glyph: "\U0000fd03", Color: [3]uint8{31, 136, 229}},  // azure
	"azure-pipelines":  {Glyph: "\U0000f427", Color: [3]uint8{20, 101, 192}},  // azure-pipelines
	"babel":            {Glyph: "\U000f0a25", Color: [3]uint8{253, 217, 59}},  // babel
	"bitbucket":        {Glyph: "\U0000f171", Color: [3]uint8{31, 136, 229}},  // bitbucket
	"blink":            {Glyph: "\U000f022b", Color: [3]uint8{249, 169, 60}},  // blink (The Foundry Nuke) (Not supported by nerdFont)
	"bower":            {Glyph: "\U0000e61a", Color: [3]uint8{239, 88, 60}},   // bower
	"bun":              {Glyph: "\U000f0685", Color: [3]uint8{249, 241, 225}}, // bun
	"c":                {Glyph: "\U000f0671", Color: [3]uint8{2, 119, 189}},   // c
	"cake":             {Glyph: "\U000f00eb", Color: [3]uint8{250, 111, 66}},  // cake
	"centOS":           {Glyph: "\U0000f304", Color: [3]uint8{157, 83, 135}},  // centOS
	"certificate":      {Glyph: "\U000f0124", Color: [3]uint8{249, 89, 63}},   // certificate
	"changelog":        {Glyph: "\U000f099b", Color: [3]uint8{139, 195, 74}},  // changelog
	"clojure":          {Glyph: "\U0000e76a", Color: [3]uint8{100, 221, 23}},  // clojure
	"cmake":            {Glyph: "\U0000f425", Color: [3]uint8{178, 178, 179}}, // cmake (Not supported by nerdFont)
	"coconut":          {Glyph: "\U000f00d3", Color: [3]uint8{141, 110, 99}},  // coconut
	"code-climate":     {Glyph: "\U000f0509", Color: [3]uint8{238, 238, 238}}, // code-climate
	"codecov":          {Glyph: "\U0000e37c", Color: [3]uint8{237, 80, 122}},  // codecov (Not supported by nerdFont)
	"codeowners":       {Glyph: "\U000f0008", Color: [3]uint8{175, 180, 43}},  // codeowners
	"coffee":           {Glyph: "\U000f0176", Color: [3]uint8{66, 165, 245}},  // coffee
	"command":          {Glyph: "\U000f0633", Color: [3]uint8{175, 188, 194}}, // command
	"commitlint":       {Glyph: "\U000f0718", Color: [3]uint8{43, 150, 137}},  // commitlint
	"conduct":          {Glyph: "\U000f014e", Color: [3]uint8{205, 220, 57}},  // conduct
	"console":          {Glyph: "\U000f018d", Color: [3]uint8{250, 111, 66}},  // console
	"contributing":     {Glyph: "\U000f014d", Color: [3]uint8{255, 202, 61}},  // contributing
	"cpp":              {Glyph: "\U000f0672", Color: [3]uint8{2, 119, 189}},   // cpp
	"credits":          {Glyph: "\U000f0260", Color: [3]uint8{156, 204, 101}}, // credits
	"csharp":           {Glyph: "\U000f031b", Color: [3]uint8{2, 119, 189}},   // csharp
	"css":              {Glyph: "\U000f031c", Color: [3]uint8{66, 165, 245}},  // css
	"css-map":          {Glyph: "\U0000e749", Color: [3]uint8{66, 165, 245}},  // css-map
	"d":                {Glyph: "\U0000e7af", Color: [3]uint8{244, 68, 62}},   // d
	"dart":             {Glyph: "\U0000e798", Color: [3]uint8{87, 182, 240}},  // dart
	"database":         {Glyph: "\U0000e706", Color: [3]uint8{255, 202, 61}},  // database
	"debian":           {Glyph: "\U0000f306", Color: [3]uint8{211, 61, 76}},   // debian
	"denizenscript":    {Glyph: "\U000f0af1", Color: [3]uint8{255, 213, 79}},  // denizenscript (Not supported by nerdFont)
	"dhall":            {Glyph: "\U0000f448", Color: [3]uint8{120, 144, 156}}, // dhall
	"disc":             {Glyph: "\U0000e271", Color: [3]uint8{176, 190, 197}}, // disc
	"django":           {Glyph: "\U0000e71d", Color: [3]uint8{67, 160, 71}},   // django
	"docker":           {Glyph: "\U0000f308", Color: [3]uint8{1, 135, 201}},   // docker
	"document":         {Glyph: "\U000f0219", Color: [3]uint8{66, 165, 245}},  // document
	"dune":             {Glyph: "\U000f0509", Color: [3]uint8{244, 127, 61}},  // dune
	"edge":             {Glyph: "\U000f0065", Color: [3]uint8{239, 111, 60}},  // edge
	"editorconfig":     {Glyph: "\U0000e652", Color: [3]uint8{255, 242, 242}}, // editorconfig
	"ejs":              {Glyph: "\U0000e618", Color: [3]uint8{255, 202, 61}},  // ejs
	"elixir":           {Glyph: "\U0000e62d", Color: [3]uint8{149, 117, 205}}, // elixir
	"elm":              {Glyph: "\U0000e62c", Color: [3]uint8{96, 181, 204}},  // elm
	"email":            {Glyph: "\U000f01ee", Color: [3]uint8{66, 165, 245}},  // email
	"erlang":           {Glyph: "\U0000e7b1", Color: [3]uint8{244, 68, 62}},   // erlang
	"eslint":           {Glyph: "\U000f0c7a", Color: [3]uint8{121, 134, 203}}, // eslint
	"exe":              {Glyph: "\U0000f2d0", Color: [3]uint8{229, 77, 58}},   // exe
	"fastlane":         {Glyph: "\U000f0700", Color: [3]uint8{149, 119, 232}}, // fastlane (Not supported by nerdFont)
	"favicon":          {Glyph: "\U0000e623", Color: [3]uint8{255, 213, 79}},  // favicon
	"fedora":           {Glyph: "\U0000f30a", Color: [3]uint8{52, 103, 172}},  // fedora
	"firebase":         {Glyph: "\U000f0967", Color: [3]uint8{251, 193, 60}},  // firebase
	"flash":            {Glyph: "\U000f0241", Color: [3]uint8{198, 52, 54}},   // flash (Not supported by nerdFont)
	"font":             {Glyph: "\U0000f031", Color: [3]uint8{244, 68, 62}},   // font
	"fortran":          {Glyph: "\U000f121a", Color: [3]uint8{250, 111, 66}},  // fortran
	"freebsd":          {Glyph: "\U0000f30c", Color: [3]uint8{175, 44, 42}},   // freebsd
	"fsharp":           {Glyph: "\U0000e7a7", Color: [3]uint8{55, 139, 186}},  // fsharp
	"gatsby":           {Glyph: "\U000f0e43", Color: [3]uint8{70, 0, 130}},    // gatsby
	"gcp":              {Glyph: "\U000f0163", Color: [3]uint8{70, 136, 250}},  // gcp (Not supported by nerdFont)
	"gemfile":          {Glyph: "\U0000e21e", Color: [3]uint8{229, 61, 58}},   // gemfile
	"gentoo":           {Glyph: "\U0000f30d", Color: [3]uint8{148, 141, 211}}, // gentoo
	"gnu":              {Glyph: "\U0000e779", Color: [3]uint8{229, 61, 58}},   // GNU
	"git":              {Glyph: "\U0000e702", Color: [3]uint8{229, 77, 58}},   // git
	"gitlab":           {Glyph: "\U0000f296", Color: [3]uint8{226, 69, 57}},   // gitlab
	"go":               {Glyph: "\U000f07d3", Color: [3]uint8{32, 173, 194}},  // go
	"go-mod":           {Glyph: "\U000f07d3", Color: [3]uint8{237, 80, 122}},  // go-mod
	"go-test":          {Glyph: "\U000f07d3", Color: [3]uint8{255, 213, 79}},  // go-test
	"godot":            {Glyph: "\U0000e65f", Color: [3]uint8{79, 195, 247}},  // godot
	"godot-assets":     {Glyph: "\U0000e65f", Color: [3]uint8{129, 199, 132}}, // godot-assets
	"gradle":           {Glyph: "\U0000e660", Color: [3]uint8{29, 151, 167}},  // gradle
	"graphql":          {Glyph: "\U000f0877", Color: [3]uint8{237, 80, 122}},  // graphql
	"groovy":           {Glyph: "\U0000f2a6", Color: [3]uint8{41, 198, 218}},  // groovy
	"grunt":            {Glyph: "\U0000e611", Color: [3]uint8{251, 170, 61}},  // grunt
	"gulp":             {Glyph: "\U0000e763", Color: [3]uint8{229, 61, 58}},   // gulp
	"h":                {Glyph: "\U000f0c00", Color: [3]uint8{2, 119, 189}},   // h (Not supported by nerdFont)
	"handlebars":       {Glyph: "\U0000e60f", Color: [3]uint8{250, 111, 66}},  // handlebars
	"haskell":          {Glyph: "\U0000e61f", Color: [3]uint8{254, 168, 62}},  // haskell
	"haxe":             {Glyph: "\U0000e666", Color: [3]uint8{246, 137, 61}},  // haxe
	"helm":             {Glyph: "\U000f0833", Color: [3]uint8{32, 173, 194}},  // helm (Not supported by nerdFont)
	"heroku":           {Glyph: "\U0000e607", Color: [3]uint8{105, 99, 185}},  // heroku
	"hpp":              {Glyph: "\U000f0b0f", Color: [3]uint8{2, 119, 189}},   // hpp (Not supported by nerdFont)
	"html":             {Glyph: "\U0000f13b", Color: [3]uint8{228, 79, 57}},   // html
	"http":             {Glyph: "\U0000f484", Color: [3]uint8{66, 165, 245}},  // http
	"husky":            {Glyph: "\U000f03e9", Color: [3]uint8{229, 229, 229}}, // husky
	"i18n":             {Glyph: "\U000f05ca", Color: [3]uint8{121, 134, 203}}, // i18n (Not supported by nerdFont)
	"image":            {Glyph: "\U000f021f", Color: [3]uint8{48, 166, 154}},  // image
	"ionic":            {Glyph: "\U0000e7a9", Color: [3]uint8{79, 143, 247}},  // ionic
	"java":             {Glyph: "\U000f0176", Color: [3]uint8{244, 68, 62}},   // java
	"javascript":       {Glyph: "\U0000e74e", Color: [3]uint8{255, 202, 61}},  // javascript
	"javascript-map":   {Glyph: "\U0000e781", Color: [3]uint8{255, 202, 61}},  // javascript-map
	"jenkins":          {Glyph: "\U0000e767", Color: [3]uint8{240, 214, 183}}, // jenkins
	"jest":             {Glyph: "\U000f0af7", Color: [3]uint8{244, 85, 62}},   // jest (Not supported by nerdFont)
	"jinja":            {Glyph: "\U0000e66f", Color: [3]uint8{174, 44, 42}},   // jinja
	"json":             {Glyph: "\U0000e60b", Color: [3]uint8{251, 193, 60}},  // json
	"julia":            {Glyph: "\U0000e624", Color: [3]uint8{134, 82, 159}},  // julia
	"karma":            {Glyph: "\U0000e622", Color: [3]uint8{60, 190, 174}},  // karma
	"key":              {Glyph: "\U000f0306", Color: [3]uint8{48, 166, 154}},  // key
	"kotlin":           {Glyph: "\U0000e634", Color: [3]uint8{139, 195, 74}},  // kotlin
	"laravel":          {Glyph: "\U0000e73f", Color: [3]uint8{248, 80, 81}},   // laravel
	"less":             {Glyph: "\U0000e60b", Color: [3]uint8{2, 119, 189}},   // less
	"lib":              {Glyph: "\U0000f831", Color: [3]uint8{139, 195, 74}},  // lib
	"linux":            {Glyph: "\U0000e712", Color: [3]uint8{238, 207, 55}},  // linux
	"liquid":           {Glyph: "\U0000e275", Color: [3]uint8{40, 182, 246}},  // liquid
	"lock":             {Glyph: "\U000f033e", Color: [3]uint8{255, 213, 79}},  // lock
	"log":              {Glyph: "\U0000f719", Color: [3]uint8{175, 180, 43}},  // log
	"lua":              {Glyph: "\U0000e620", Color: [3]uint8{66, 165, 245}},  // lua
	"makefile":         {Glyph: "\U000f0229", Color: [3]uint8{239, 83, 80}},   // makefile
	"manjaro":          {Glyph: "\U0000f312", Color: [3]uint8{73, 185, 90}},   // manjaro
	"markdown":         {Glyph: "\U0000eb1d", Color: [3]uint8{66, 165, 245}},  // markdown
	"markojs":          {Glyph: "\U0000f13b", Color: [3]uint8{2, 119, 189}},   // markojs (Not supported by nerdFont)
	"mdx":              {Glyph: "\U0000f853", Color: [3]uint8{255, 202, 61}},  // mdx
	"merlin":           {Glyph: "\U0000f136", Color: [3]uint8{66, 165, 245}},  // merlin (Not supported by nerdFont)
	"midi":             {Glyph: "\U000f08f2", Color: [3]uint8{66, 165, 245}},  // midi
	"mint":             {Glyph: "\U0000f30f", Color: [3]uint8{125, 190, 58}},  // mint
	"mjml":             {Glyph: "\U0000e714", Color: [3]uint8{249, 89, 63}},   // mjml (Not supported by nerdFont)
	"mocha":            {Glyph: "\U000f01aa", Color: [3]uint8{161, 136, 127}}, // mocha (Not supported by nerdFont)
	"modernizr":        {Glyph: "\U0000e720", Color: [3]uint8{234, 72, 99}},   // modernizr
	"moonscript":       {Glyph: "\U0000f186", Color: [3]uint8{251, 193, 60}},  // moonscript
	"mxml":             {Glyph: "\U0000f72d", Color: [3]uint8{254, 168, 62}},  // mxml
	"mysql":            {Glyph: "\U0000e704", Color: [3]uint8{1, 94, 134}},    // mysql
	"nim":              {Glyph: "\U000f01a5", Color: [3]uint8{255, 202, 61}},  // nim
	"nix":              {Glyph: "\U0000f313", Color: [3]uint8{80, 117, 193}},  // nix
	"nodejs":           {Glyph: "\U000f0399", Color: [3]uint8{139, 195, 74}},  // nodejs
	"npm":              {Glyph: "\U0000e71e", Color: [3]uint8{203, 56, 55}},   // npm
	"nuget":            {Glyph: "\U0000e77f", Color: [3]uint8{3, 136, 209}},   // nuget (Not supported by nerdFont)
	"nuxt":             {Glyph: "\U000f1106", Color: [3]uint8{65, 184, 131}},  // nuxt
	"ocaml":            {Glyph: "\U0000e67a", Color: [3]uint8{253, 154, 62}},  // ocaml
	"opam":             {Glyph: "\U0000f1ce", Color: [3]uint8{255, 213, 79}},  // opam (Not supported by nerdFont)
	"opensuse":         {Glyph: "\U0000f314", Color: [3]uint8{111, 180, 36}},  // opensuse
	"pascal":           {Glyph: "\U000f03db", Color: [3]uint8{3, 136, 209}},   // pascal (Not supported by nerdFont)
	"pawn":             {Glyph: "\U0000e261", Color: [3]uint8{239, 111, 60}},  // pawn
	"pdf":              {Glyph: "\U0000e67d", Color: [3]uint8{244, 68, 62}},   // pdf
	"perl":             {Glyph: "\U0000e769", Color: [3]uint8{149, 117, 205}}, // perl
	"php":              {Glyph: "\U0000e608", Color: [3]uint8{65, 129, 190}},  // php
	"postcss":          {Glyph: "\U000f031c", Color: [3]uint8{244, 68, 62}},   // postcss (Not supported by nerdFont)
	"postgresql":       {Glyph: "\U0000e76e", Color: [3]uint8{49, 99, 140}},   // postgresql
	"powerpoint":       {Glyph: "\U000f0227", Color: [3]uint8{209, 71, 51}},   // powerpoint
	"powershell":       {Glyph: "\U000f0a0a", Color: [3]uint8{5, 169, 244}},   // powershell
	"prettier":         {Glyph: "\U0000e6b4", Color: [3]uint8{86, 179, 180}},  // prettier
	"prisma":           {Glyph: "\U0000e684", Color: [3]uint8{229, 56, 150}},  // prisma
	"prolog":           {Glyph: "\U0000e7a1", Color: [3]uint8{239, 83, 80}},   // prolog
	"protractor":       {Glyph: "\U0000f288", Color: [3]uint8{229, 61, 58}},   // protractor (Not supported by nerdFont)
	"pug":              {Glyph: "\U0000e686", Color: [3]uint8{239, 204, 163}}, // pug
	"puppet":           {Glyph: "\U0000e631", Color: [3]uint8{251, 193, 60}},  // puppet
	"purescript":       {Glyph: "\U0000e630", Color: [3]uint8{66, 165, 245}},  // purescript
	"python":           {Glyph: "\U000f0320", Color: [3]uint8{52, 102, 143}},  // python
	"python-misc":      {Glyph: "\U0000f820", Color: [3]uint8{130, 61, 28}},   // python-misc
	"qsharp":           {Glyph: "\U0000f292", Color: [3]uint8{251, 193, 60}},  // qsharp (Not supported by nerdFont)
	"r":                {Glyph: "\U000f07d4", Color: [3]uint8{25, 118, 210}},  // r
	"raml":             {Glyph: "\U0000e60b", Color: [3]uint8{66, 165, 245}},  // raml
	"raspberry-pi":     {Glyph: "\U0000f315", Color: [3]uint8{208, 60, 76}},   // raspberry-pi
	"razor":            {Glyph: "\U000f0065", Color: [3]uint8{66, 165, 245}},  // razor
	"react":            {Glyph: "\U0000e7ba", Color: [3]uint8{35, 188, 212}},  // react
	"react_ts":         {Glyph: "\U0000e7ba", Color: [3]uint8{36, 142, 211}},  // react_ts
	"readme":           {Glyph: "\U000f02fc", Color: [3]uint8{66, 165, 245}},  // readme
	"redhat":           {Glyph: "\U0000f316", Color: [3]uint8{231, 61, 58}},   // redhat
	"roadmap":          {Glyph: "\U000f066e", Color: [3]uint8{48, 166, 154}},  // roadmap
	"robot":            {Glyph: "\U000f06a9", Color: [3]uint8{249, 89, 63}},   // robot
	"rollup":           {Glyph: "\U000f0bc0", Color: [3]uint8{233, 61, 50}},   // rollup
	"routing":          {Glyph: "\U000f0641", Color: [3]uint8{67, 160, 71}},   // routing
	"ruby":             {Glyph: "\U0000e739", Color: [3]uint8{229, 61, 58}},   // ruby
	"rust":             {Glyph: "\U0000e7a8", Color: [3]uint8{250, 111, 66}},  // rust
	"sass":             {Glyph: "\U0000e603", Color: [3]uint8{237, 80, 122}},  // sass
	"scheme":           {Glyph: "\U000f0627", Color: [3]uint8{244, 68, 62}},   // scheme
	"semantic-release": {Glyph: "\U000f0210", Color: [3]uint8{245, 245, 245}}, // semantic-release (Not supported by nerdFont)
	"settings":         {Glyph: "\U0000f013", Color: [3]uint8{66, 165, 245}},  // settings
	"shaderlab":        {Glyph: "\U000f06af", Color: [3]uint8{25, 118, 210}},  // shaderlab
	"sketch":           {Glyph: "\U000f0b8a", Color: [3]uint8{255, 194, 61}},  // sketch
	"slim":             {Glyph: "\U0000f24e", Color: [3]uint8{245, 129, 61}},  // slim (Not supported by nerdFont)
	"smarty":           {Glyph: "\U000f0335", Color: [3]uint8{255, 207, 60}},  // smarty
	"solidity":         {Glyph: "\U000f07bb", Color: [3]uint8{3, 136, 209}},   // solidity
	"sqlite":           {Glyph: "\U0000e7c4", Color: [3]uint8{1, 57, 84}},     // sqlite
	"storybook":        {Glyph: "\U000f082e", Color: [3]uint8{237, 80, 122}},  // storybook (Not supported by nerdFont)
	"stryker":          {Glyph: "\U0000f05b", Color: [3]uint8{239, 83, 80}},   // stryker
	"stylelint":        {Glyph: "\U0000e695", Color: [3]uint8{207, 216, 220}}, // stylelint
	"stylus":           {Glyph: "\U0000e600", Color: [3]uint8{192, 202, 51}},  // stylus
	"sublime":          {Glyph: "\U0000e696", Color: [3]uint8{239, 148, 58}},  // sublime
	"svelte":           {Glyph: "\U0000e697", Color: [3]uint8{255, 62, 0}},    // svelte
	"svg":              {Glyph: "\U000f0721", Color: [3]uint8{255, 181, 62}},  // svg
	"swc":              {Glyph: "\U000f06d5", Color: [3]uint8{198, 52, 54}},   // swc (Not supported by nerdFont)
	"swift":            {Glyph: "\U000f06e5", Color: [3]uint8{249, 95, 63}},   // swift
	"table":            {Glyph: "\U000f021b", Color: [3]uint8{139, 195, 74}},  // table
	"tailwindcss":      {Glyph: "\U000f13ff", Color: [3]uint8{77, 182, 172}},  // tailwindcss
	"tcl":              {Glyph: "\U000f06d3", Color: [3]uint8{239, 83, 80}},   // tcl
	"terraform":        {Glyph: "\U000f1062", Color: [3]uint8{92, 107, 192}},  // terraform
	"test-js":          {Glyph: "\U000f0096", Color: [3]uint8{255, 202, 61}},  // test-js
	"test-jsx":         {Glyph: "\U000f0096", Color: [3]uint8{35, 188, 212}},  // test-jsx
	"test-ts":          {Glyph: "\U000f0096", Color: [3]uint8{3, 136, 209}},   // test-ts
	"tex":              {Glyph: "\U0000e69b", Color: [3]uint8{66, 165, 245}},  // tex
	"todo":             {Glyph: "\U0000f058", Color: [3]uint8{124, 179, 66}},  // todo
	"travis":           {Glyph: "\U0000e77e", Color: [3]uint8{203, 58, 73}},   // travis
	"tune":             {Glyph: "\U000f066a", Color: [3]uint8{251, 193, 60}},  // tune
	"twig":             {Glyph: "\U0000e61c", Color: [3]uint8{155, 185, 47}},  // twig
	"typescript":       {Glyph: "\U0000e628", Color: [3]uint8{3, 136, 209}},   // typescript
	"typescript-def":   {Glyph: "\U000f06e6", Color: [3]uint8{3, 136, 209}},   // typescript-def
	"typst":            {Glyph: "t", Color: [3]uint8{34, 157, 172}},           // typst
	"ubuntu":           {Glyph: "\U0000f31c", Color: [3]uint8{214, 73, 53}},   // ubuntu
	"url":              {Glyph: "\U000f0337", Color: [3]uint8{66, 165, 245}},  // url
	"vagrant":          {Glyph: "\U0000f27d", Color: [3]uint8{20, 101, 192}},  // vagrant (Not supported by nerdFont)
	"vala":             {Glyph: "\U0000e69e", Color: [3]uint8{149, 117, 205}}, // vala (Not supported by nerdFont)
	"vercel":           {Glyph: "\U0000f47e", Color: [3]uint8{207, 216, 220}}, // vercel
	"verilog":          {Glyph: "\U000f061a", Color: [3]uint8{250, 111, 66}},  // verilog
	"video":            {Glyph: "\U000f022b", Color: [3]uint8{253, 154, 62}},  // video
	"vim":              {Glyph: "\U0000e62b", Color: [3]uint8{67, 160, 71}},   // vim
	"virtual":          {Glyph: "\U0000f822", Color: [3]uint8{3, 155, 229}},   // virtual
	"visualstudio":     {Glyph: "\U0000e70c", Color: [3]uint8{173, 99, 188}},  // visualstudio
	"vscode":           {Glyph: "\U0000e70c", Color: [3]uint8{33, 150, 243}},  // vscode (Not supported by nerdFont)
	"vue":              {Glyph: "\U000f0844", Color: [3]uint8{65, 184, 131}},  // vue
	"vue-config":       {Glyph: "\U000f0844", Color: [3]uint8{58, 121, 110}},  // vue-config
	"webpack":          {Glyph: "\U000f072b", Color: [3]uint8{142, 214, 251}}, // webpack
	"word":             {Glyph: "\U000f022c", Color: [3]uint8{1, 87, 155}},    // word
	"xaml":             {Glyph: "\U000f05c0", Color: [3]uint8{64, 153, 69}},   // xaml
	"xml":              {Glyph: "\U000f05c0", Color: [3]uint8{64, 153, 69}},   // xml
	"yaml":             {Glyph: "\U0000e60b", Color: [3]uint8{244, 68, 62}},   // yaml
	"yang":             {Glyph: "\U000f0680", Color: [3]uint8{66, 165, 245}},  // yang
	"yarn":             {Glyph: "\U000f011b", Color: [3]uint8{44, 142, 187}},  // yarn
	"zig":              {Glyph: "\U0000e6a9", Color: [3]uint8{249, 169, 60}},  // zig
	"zip":              {Glyph: "\U0000f410", Color: [3]uint8{175, 180, 43}},  // zip

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

	"dir-config":      {Glyph: "\U0000e5fc", Color: [3]uint8{32, 173, 194}},  // dir-config
	"dir-controller":  {Glyph: "\U0000e5fc", Color: [3]uint8{255, 194, 61}},  // dir-controller
	"dir-download":    {Glyph: "\U000f024d", Color: [3]uint8{76, 175, 80}},   // dir-download
	"dir-environment": {Glyph: "\U000f024f", Color: [3]uint8{102, 187, 106}}, // dir-environment
	"dir-git":         {Glyph: "\U0000e5fb", Color: [3]uint8{250, 111, 66}},  // dir-git
	"dir-github":      {Glyph: "\U0000e5fd", Color: [3]uint8{84, 110, 122}},  // dir-github
	"dir-images":      {Glyph: "\U000f024f", Color: [3]uint8{43, 150, 137}},  // dir-images
	"dir-import":      {Glyph: "\U000f0257", Color: [3]uint8{175, 180, 43}},  // dir-import
	"dir-include":     {Glyph: "\U000f0257", Color: [3]uint8{3, 155, 229}},   // dir-include
	"dir-npm":         {Glyph: "\U0000e5fa", Color: [3]uint8{203, 56, 55}},   // dir-npm
	"dir-secure":      {Glyph: "\U000f0250", Color: [3]uint8{249, 169, 60}},  // dir-secure
	"dir-upload":      {Glyph: "\U000f0259", Color: [3]uint8{250, 111, 66}},  // dir-upload
}

// default icons in case nothing can be found
var IconDef = map[string]*IconInfo{
	"dir":        {Glyph: "\U000f024b", Color: [3]uint8{224, 177, 77}},
	"diropen":    {Glyph: "\U000f0770", Color: [3]uint8{224, 177, 77}},
	"exe":        {Glyph: "\U000f0214", Color: [3]uint8{76, 175, 80}},
	"file":       {Glyph: "\U000f0224", Color: [3]uint8{65, 129, 190}},
	"hiddendir":  {Glyph: "\U000f0256", Color: [3]uint8{224, 177, 77}},
	"hiddenfile": {Glyph: "\U000f0613", Color: [3]uint8{65, 129, 190}},
}
