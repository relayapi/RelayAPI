openai 这类的 api 调用，需要使用api key，但是 key 不能直接暴露在前端，只能在后端使用。
前端需要调用后端的接口，后端再调用 openai 的 api，中间的过程很繁琐，需要处理很多问题。

我想实现一个中间层，前端调用中间层的接口，中间层调用 openai 的 api。

这个中间层名字叫 RelayAPI, 它是一个代理层，可以代理任何 api 调用，只需要在 RelayAPI 中配置好公钥，私钥由后端保存，后端把请求的 明文参数（模型名称、APIkey）用私钥加密，生成密文参数，把URL(RelayAPI服务器地址)+参数交给前端，前端调用 RelayAPI 的接口，RelayAPI 再把密文参数用公钥解密，得到明文参数，再调用 openai 的 api。

RelayAPI 包括两部分，一部分是RelayAPI server，另一部分是RelayAPI client SDK ，后端使用client SDK 生成参数，交给前端，再由前端调用 RelayAPI server 的接口，RelayAPI server 再把密文参数用公钥解密，得到明文参数，再调用 openai 的 api。

为了避免前端拿到密文参数后无限制的调用，在参数中可以设置调用次数和调用时间，后端可以随时查看调用次数和调用时间，如果调用次数和调用时间超限，则拒绝调用。

RelayAPI server 使用 Golang 开发，支持高速、高并发,非对称加密使用ecc，支持多语言的 client SDK。

