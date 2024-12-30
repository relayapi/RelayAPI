import axios from 'axios';
import { TokenGenerator } from './TokenGenerator.js';

export class RelayAPIClient {
    /**
     * Initialize RelayAPI client
     * 初始化 RelayAPI 客户端
     * @param {object} config Configuration object / 配置对象
     */
    constructor(config) {
        this.tokenGenerator = new TokenGenerator(config);
    }

    /**
     * Create token
     * 创建令牌
     * @param {object} options Token options / 令牌选项
     * @returns {string} Encrypted token string / 加密的令牌字符串
     */
    createToken(options) {
        const tokenData = this.tokenGenerator.createToken(options);
        return this.tokenGenerator.encryptToken(tokenData);
    }

    /**
     * Generate API URL
     * 生成 API URL
     * @param {string} token Token / 令牌
     * @param {string} [endpoint=''] API endpoint / API 端点
     * @returns {string} Complete API URL / 完整的 API URL
     */
    generateUrl(token, endpoint = '') {
        const baseUrl = this.tokenGenerator.getServerUrl(endpoint);
        return `${baseUrl}?token=${token}&rai_hash=${this.tokenGenerator.hash}`;
    }

    /**
     * Send chat request
     * 发送聊天请求
     * @param {object} options Request options / 请求选项
     * @returns {Promise<object>} Response data / 响应数据
     */
    async chat(options) {
        const {
            messages,
            model = 'gpt-3.5-turbo',
            temperature = 0.7,
            maxTokens = 1000,
            token
        } = options;

        const url = this.generateUrl(token, 'chat/completions');
        
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
     * Generate image
     * 生成图像
     * @param {object} options Request options / 请求选项
     * @returns {Promise<object>} Response data / 响应数据
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

        const url = this.generateUrl(token, 'images/generations');
        
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
     * Generate embeddings
     * 生成嵌入向量
     * @param {object} options Request options / 请求选项
     * @returns {Promise<object>} Response data / 响应数据
     */
    async createEmbedding(options) {
        const {
            input,
            model = 'text-embedding-ada-002',
            token
        } = options;

        const url = this.generateUrl(token, 'embeddings');
        
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
     * Check server health status
     * 检查服务器健康状态
     * @returns {Promise<object>} Health status data / 健康状态数据
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