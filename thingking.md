# aic开发提示词

开发命令行工具，连接本地的ollama接口，执行`aic <prompt>`，然后转换成命令行的形式执行，把运行的结果输出。

## 提示词要求

- 模型默认选择`qwen2.5-coder`，支持`model`参数修改
- 系统提示词必须严格约束只根据用户的描述生成当前系统环境的命令行
  - 支持darwin，linux，windows
  - 支持`zsh`，`bash`，`powershell`，`cmd`
  - 支持不同的操作系统版本
  - 系统提示词中增加当前操作系统和shell的具体信息
  - 系统提示词中增加当前用户的用户名和用户目录，以及当前所在目录
  - 如果无法生成命令，应当放回`<err_cannot_generate_command>`错误，方便提示用户
- 提示词要考虑命令行的正确性
  - 比如生成的命令行是`curl -s "http://a.com/?q=%C %t"`会报错`curl: (3) URL rejected: Malformed input to a URL function`，因为url格式不对，需要避免这种情况，考虑url编码
  - 尽可能避开需要api_key的场景，如果确实需要，先确保环境变量中有相关的配置
  - 涉及到互联网请求，有限判断如下条件：
    - 如果是询问天气，访问`https://wttr.in/`

## 程序交互要求

- ollama的地址可以通过参数`ollama-url`指定，默认是`http://localhost:11434`
- 支持调试模式`verbose`参数
  - 打印版本号
  - 打印调用的ollama地址和模型
  - 打印用户输入的prompt
  - 打印ollama返回的实际命令行
- 模型的返回要考虑到对错误的处理
  - 模型未找到的情况
  
## 基本框架要求

- 用golang作为开发语言，目录结构符合golang最佳实践
- 代码模块化
- 有完善的单元测试，确保代码覆盖度尽可能高
- 用golangci-lint检查代码，满足最佳实践，修复代码中的警告和错误
- 完善`README.md`文档，遵循知名开源项目的最佳实践，包括说明，Feature，安装，运行，以及各种徽章等
  - 徽章要求质量必须是最好的，至少支持：
    - `https://goreportcard.com`
- 支持goreleaser跨平台编译
  - goreleaser在打包过程中必须要是编译出用于发布的最小化程序
- 支持github action自动执行测试和编译
- 支持vscode调试（.vscode/launch.json），配置文件中的参数为`-verbose hi`
- 命令行的返回要美观大方，包括支持颜色，在错误的情况下用醒目的颜色输出
- git提交排除`dist`这种临时或者生成的目录
- 程序里面的所有给用户的输出都需要是中文，注释可以是英文
- 生成一个makefile
  - 默认是build编译
  - 支持test完成测试和覆盖度的输出
  - 支持release完成测试编译以及github的release文件上传
    - make release的过程中要提示用户输入tag，否则就带上snapshot参数，否则会提示错误：`error=git doesn't contain any tags. Either add a tag or use --snapshot`
    - 提交前避免出现`error=git is in a dirty state`这样的错误，就先检查一下，并且让用户确认是否把当前未提交的代码先commit再自动push提交，确保后续能够正确的推送github release包
    - 提交前避免出现类似`error=git tag v0.0.1 was not made against commit a95e41ef507e852569a36bb2c5377b325a7e24a2`这样的错误，如果已经有了tag，就让用户确认是新建还是删除再重新创建相同的，如果是新建就让用户输入新的tag值。重新创建的时候也要把服务器上的tag也删除，同时把上面所有相同tag的release的draft全部删除。
  - 新增CHANGELOG.md文件来记录版本号和对应的更新情况，在release的时候自动读取并且生成对应的版本号推送，在程序里面可以通过`aic -verion`来查看
  - 除了代码注释和md文档可以出现中文外，所有命令行输出、打包、测试、提交、推送、release等让用户看得到的所有操作都必须是英文