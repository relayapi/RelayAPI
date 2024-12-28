import axios from 'axios';
import { TokenGenerator } from './TokenGenerator.js';

export class RelayAPIClient {
    /**
     * 初始化 RelayAPI 客户端
     * @param {string|object} config 配置文件路径或配置对象
     */
    constructor(config) {
        this.tokenGenerator = new TokenGenerator(config);
    }

    /**
     * 创建令牌
     * @param {object} options 令牌选项
     * @returns {string} 加密的令牌字符串
     */
    createToken(options) {
        const tokenData = this.tokenGenerator.createToken(options);
        return this.tokenGenerator.encryptToken(tokenData);
    }

    /**
     * 生成 API URL
     * @param {string} endpoint API 端点
     * @param {string} token 令牌
     * @returns {string} 完整的 API URL
     */
    generateUrl(endpoint, token) {
        const baseUrl = this.tokenGenerator.getServerUrl(endpoint);
        return `${baseUrl}?token=${token}`;
    }

    /**
     * 发送聊天请求
     * @param {object} options 请求选项
     * @returns {Promise<object>} 响应数据
     */
    async chat(options) {
        const {
            messages,
            model = 'gpt-3.5-turbo',
            temperature = 0.7,
            maxTokens = 1000,
            token
        } = options;

        const url = this.generateUrl('v1/chat/completions', token);
        
        try {
            const response = await axios.post(url, {
                messages,
                model,
                temperature,
                max_tokens: maxTokens
            });
            return response.data;
        } catch (error) {
            throw new Error(`Chat request failed: ${error.message}`);
        }
    }

    /**
     * 生成图像
     * @param {object} options 请求选项
     * @returns {Promise<object>} 响应数据
     */
    async generateImage(options) {
        const {
            prompt,
            model = 'dall-e-3',
            size = '1024x1024',
            quality = 'standard',
            n = 1,
            token
        } = options;

        const url = this.generateUrl('v1/images/generations', token);
        
        try {
            const response = await axios.post(url, {
                prompt,
                model,
                size,
                quality,
                n
            });
            return response.data;
        } catch (error) {
            throw new Error(`Image generation failed: ${error.message}`);
        }
    }

    /**
     * 生成嵌入向量
     * @param {object} options 请求选项
     * @returns {Promise<object>} 响应数据
     */
    async createEmbedding(options) {
        const {
            input,
            model = 'text-embedding-ada-002',
            token
        } = options;

        const url = this.generateUrl('v1/embeddings', token);
        
        try {
            const response = await axios.post(url, {
                input,
                model
            });
            return response.data;
        } catch (error) {
            throw new Error(`Embedding creation failed: ${error.message}`);
        }
    }

    /**
     * 检查服务器健康状态
     * @returns {Promise<object>} 健康状态数据
     */
    async healthCheck() {
        const url = this.tokenGenerator.getServerUrl('health');
        try {
            const response = await axios.get(url);
            return response.data;
        } catch (error) {
            throw new Error(`Health check failed: ${error.message}`);
        }
    }
} 