<div align=center>

<img src="https://iim.ai/public/images/logo.png?t=20231207" width="128"/>

# 智元 Fast API SDK

[![Release](https://img.shields.io/github/v/release/iimeta/fastapi-sdk?color=blue)](https://github.com/iimeta/fastapi-sdk/releases)
[![License](https://img.shields.io/static/v1?label=license&message=MIT&color=green)](https://github.com/iimeta/fastapi-sdk/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/iimeta/fastapi-sdk)](https://goreportcard.com/report/github.com/iimeta/fastapi-sdk)

</div>

`智元 Fast API` 是一站式API管理系统，将各类LLM API进行统一格式、统一规范、统一管理，使其在功能、性能和用户体验上达到极致。

## 功能介绍
智元 Fast API SDK模块, 主要集成各类 `LLM API` 的能力

## 支持模型

| 公司                                    | 模型                                                                                                                                                                                     | Completion | Image | Audio | Embedding | Multimodal |
|-----------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|------------|-------|-------|-----------|------------|
| OpenAI<br/>Azure                        | gpt3.5、4.0系列、chatgpt-4o-latest、o1-mini、o1-preview<br/>DALL-E-2、DALL-E-3<br/>tts-1、tts-1-hd、whisper-1<br/>text-embedding-3-small、text-embedding-3-large、text-embedding-ada-002 | ✔️         | ✔️    | ✔️    | ✔️        | ✔️         |
| 百度                                    | 文心一言                                                                                                                                                                                 | ✔️         |       |       |           |            |
| 科大讯飞                                | 星火                                                                                                                                                                                     | ✔️         | ✔️    |       |           |            |
| 阿里云                                  | 通义千问                                                                                                                                                                                 | ✔️         |       |       |           |            |
| 智谱AI                                  | GLM系列                                                                                                                                                                                  | ✔️         |       |       |           |            |
| Google                                  | Gemini                                                                                                                                                                                   | ✔️         |       |       |           | ✔️         |
| DeepSeek                                | deepseek-chat、deepseek-coder                                                                                                                                                            | ✔️         |       |       |           |            |
| 360智脑                                 | 360gpt系列                                                                                                                                                                               | ✔️         |       |       |           |            |
| Midjourney                              | Midjourney-Proxy、Midjourney-Proxy-Plus                                                                                                                                                  |            | ✔️    |       |           |            |
| Anthropic<br/>GCPClaude<br/>AWSClaude   | claude-3-5-sonnet-20240620<br/>claude-3-opus-20240229<br/>claude-3-haiku-20240307                                                                                                        | ✔️         |       |       |           |            |

## 项目预览

- 用户端: https://demo.fastapi.ai/login
- 账号/密码: demo@fastapi.ai/123456

- 管理端: https://demo.fastapi.ai/admin
- 账号/密码: admin/admin123

## 公益站点

- API地址: https://api.free.fastapi.ai
- 注册后请联系作者领取1000万额度

## 项目截图
### 管理端
![dashboard](https://github.com/iimeta/fastapi-admin/assets/138393700/4d26474f-2082-41b4-8bd9-843ccb66d0a7)
![model](https://github.com/iimeta/fastapi-admin/assets/138393700/89d0721c-a72b-47f3-84e3-0080f5f1c8f9)
![key](https://github.com/iimeta/fastapi-admin/assets/138393700/f93922de-92e2-438e-8786-ae26981ec154)
![modelagent](https://github.com/iimeta/fastapi-admin/assets/138393700/f1e2c21f-d5ce-45e7-811d-7a5ed1a04138)
![user](https://github.com/iimeta/fastapi-admin/assets/138393700/1e812ed7-ff9f-46ab-847d-23a177bbedf0)
![chat](https://github.com/iimeta/fastapi-admin/assets/138393700/041a4324-eba8-42d4-8467-5964df6242f3)

### 用户端
![dashboard](https://github.com/iimeta/fastapi-admin/assets/138393700/94cd152c-b370-4cd9-b3f2-ca1dddf6821e)
![mymodel](https://github.com/iimeta/fastapi-admin/assets/138393700/cd8fec82-b2b8-4af4-b471-2c313e321d30)
![app](https://github.com/iimeta/fastapi-admin/assets/138393700/4dec9a5f-f399-4bfc-a3e3-bb50248a2c4c)
![appkey](https://github.com/iimeta/fastapi-admin/assets/138393700/662fe51d-9ed1-4672-b9db-f418b492b9a2)
![chat](https://github.com/iimeta/fastapi-admin/assets/138393700/5a92ab02-1319-436c-af2d-7d9ef0f49fe5)

## 开源地址
| 仓库 | API                                          | Web                                                  | Admin                                                    | SDK                                                  |
| ---- |----------------------------------------------|------------------------------------------------------|----------------------------------------------------------|------------------------------------------------------|
| 主库 | [fastapi](https://github.com/iimeta/fastapi) | [fastapi-web](https://github.com/iimeta/fastapi-web) | [fastapi-admin](https://github.com/iimeta/fastapi-admin) | [fastapi-sdk](https://github.com/iimeta/fastapi-sdk) |
| 码云 | [fastapi](https://gitee.com/iimeta/fastapi)  | [fastapi-web](https://gitee.com/iimeta/fastapi-web)  | [fastapi-admin](https://gitee.com/iimeta/fastapi-admin)  | [fastapi-sdk](https://gitee.com/iimeta/fastapi-sdk)  |

## 作者
<img src="https://iim.ai/public/images/Author.png?t=20231207" width="320" alt="作者"/>

## 交流群
<img src="https://iim.ai/public/images/WeChatGroup.jpg?t=20241001" width="320" alt="交流群"/>


![Star History Chart](https://api.star-history.com/svg?repos=iimeta/fastapi-sdk&type=Date)
