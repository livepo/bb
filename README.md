# bb

轻量级终端翻译工具，基于百度翻译 API。适用于在命令行下快速翻译短句或单词，支持中英互译并可通过 Homebrew 安装。

## 特性

- 极简的命令行界面
- 支持自动检测源语言（中/英互转）
- 可通过 Homebrew Tap 安装（如果已发布）或使用提供的安装脚本安装

## 安装

推荐方法（Homebrew）：

1. 添加 tap（示例）：

	```bash
	brew tap livepo/homebrew-tap
	```

2. 安装：

	```bash
	brew install bb
	```

手动安装（使用安装脚本）：

```bash
curl -sL https://github.com/livepo/bb/raw/master/install.sh | bash
```

或者直接下载并复制到 `/usr/local/bin`：

```bash
wget -O /usr/local/bin/bb https://github.com/livepo/bb/releases/download/vX.Y.Z/bb
chmod +x /usr/local/bin/bb
```

> 注：替换上面的 `vX.Y.Z` 为实际发布的版本号。

## 配置

在首次使用前，请在你的家目录中创建配置文件 `~/.bb`，并包含以下三项（分别为 APPID、SECRET、SALT）：

```
APPID=你的_appid
SECRET=你的_secret
SALT=一个随机字符串
```

如果缺少任一配置项，程序会退出并提示错误。

## 使用方法

基本命令：

```bash
bb <要翻译的文本>
```

示例：

```bash
bb 你好
bb Hello
```

查看版本：

```bash
bb version
```

## 输出示例

命令：

```bash
bb Hello
```

示例输出：

```
Hello
-------------
你好
```

## 贡献

欢迎 PR、Issue 和建议。贡献前请确保：

- 遵循仓库已有的代码风格
- 添加必要的测试（如果适用）

## 许可证

此项目采用 MIT 许可证。详情见 `LICENSE`
