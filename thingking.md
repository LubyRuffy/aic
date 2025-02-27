# aic开发提示词

开发命令行工具，连接本地的ollama接口，执行`aic <prompt>`，然后转换成命令行的形式执行，把运行的结果输出。

要求：

- 用golang作为开发语言，目录结构符合golang最佳实践
- 代码模块化
- 有完善的单元测试，确保代码覆盖度尽可能高
- 用golangci-lint检查代码，满足最佳实践，修复代码中的警告和错误
- 模型默认选择`qwen2.5-coder`，支持`model`参数修改
- 模型的返回要考虑到对错误的处理
  - 模型未找到的情况
- 系统提示词必须严格约束只根据用户的描述生成当前系统环境的命令行
  - 支持darwin，linux，windows
  - 支持`zsh`，`bash`，`powershell`，`cmd`
  - 支持不同的操作系统版本
  - 系统提示词中增加当前操作系统和shell的具体信息
  - 系统提示词中增加当前用户的用户名和用户目录，以及当前所在目录
  - 如果无法生成命令，应当放回`<err_cannot_generate_command>`错误，方便提示用户
- ollama的地址可以通过参数`ollama-url`指定，默认是`http://localhost:11434`
- 支持调试模式`verbose`参数，可以打印ollama返回的实际命令行
- 完善`README.md`文档，遵循知名开源项目的最佳实践，包括说明，Feature，安装，运行，以及各种徽章等
- 支持goreleaser跨平台编译
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