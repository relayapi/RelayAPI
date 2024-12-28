import { RelayAPIClient } from '../src/Client.js';
import fs from 'fs/promises';

async function main() {
    try {
        // 从配置文件加载配置
        const configContent = await fs.readFile('default.rai', 'utf-8');
        const config = JSON.parse(configContent);

        // 创建客户端实例
        const client = new RelayAPIClient(config);

        // 创建令牌
        const token = client.createToken({
            apiKey: 'sk-573af3eca24f492a83d5e64894ed91f5',
            maxCalls: 100,
            expireDays: 1,
            provider: 'dashscope'
        });

        console.log('Token created:', token);

        // 检查服务器健康状态
        const healthStatus = await client.healthCheck();
        console.log('Server health status:', healthStatus);

        // 发送聊天请求
        const chatResponse = await client.chat({
            messages: [
                { role: 'system', content: 'You are a helpful assistant.' },
                { role: 'user', content: 'What is the capital of France?' }
            ],
            model: 'qwen-vl-max',
            temperature: 0.7,
            maxTokens: 1000,
            token: token
        });

        console.log('Chat response:', chatResponse);

        // 生成图像
        const imageResponse = await client.generateImage({
            prompt: 'A beautiful sunset over Paris',
            model: 'dall-e-3',
            size: '1024x1024',
            quality: 'standard',
            n: 1,
            token: token
        });

        console.log('Image response:', imageResponse);

        // 生成嵌入向量
        const embeddingResponse = await client.createEmbedding({
            input: 'The quick brown fox jumps over the lazy dog',
            model: 'text-embedding-ada-002',
            token: token
        });

        console.log('Embedding response:', embeddingResponse);

    } catch (error) {
        console.error('Error:', error.message);
    }
}

// 运行示例程序
main(); 