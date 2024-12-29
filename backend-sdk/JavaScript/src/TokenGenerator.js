import CryptoJS from 'crypto-js';
import fs from 'fs/promises';

export class TokenGenerator {
    /**
     * Initialize token generator
     * 初始化令牌生成器
     * @param {string|object} config Configuration file path or configuration object / 配置文件路径或配置对象
     */
    constructor(config) {
        this.config = null;
        this.initialize(config);
    }

    /**
     * Initialize configuration
     * 初始化配置
     * @param {string|object} config Configuration file path or configuration object / 配置文件路径或配置对象
     */
    async initialize(config) {
        this.config = config;

        // Validate configuration / 验证配置
        this.validateConfig();
    }

    /**
     * Validate configuration
     * 验证配置
     */
    validateConfig() {
        const { crypto, server } = this.config;
        if (!crypto || !server) {
            throw new Error('Invalid config: missing crypto or server section');
        }
        if (crypto.method !== 'aes') {
            throw new Error('Only AES encryption is supported');
        }
        if (!crypto.aes_key || !crypto.aes_iv_seed) {
            throw new Error('Invalid crypto config: missing aes_key or aes_iv_seed');
        }
    }

    /**
     * Generate initialization vector
     * 生成初始化向量
     * @returns {WordArray} Generated IV / 生成的 IV
     */
    generateIV() {
        // Generate random IV / 生成随机 IV
        const iv = CryptoJS.lib.WordArray.random(16);
        // Mix with IV seed / 使用 IV 种子进行混合
        const ivSeed = CryptoJS.enc.Utf8.parse(this.config.crypto.aes_iv_seed);
        const words = iv.words.map((word, i) => word ^ ivSeed.words[i]);
        return CryptoJS.lib.WordArray.create(words, 16);
    }

    /**
     * Create token data
     * 创建令牌数据
     * @param {object} options Token options / 令牌选项
     * @returns {object} Token data / 令牌数据
     */
    createToken({
        apiKey,
        maxCalls = 100,
        expireDays = 1,
        provider = 'dashscope',
        extInfo = ''
    }) {
        const now = new Date();
        const expireTime = new Date(now.getTime() + expireDays * 24 * 60 * 60 * 1000);

        return {
            id: `token-${Date.now()}`,
            api_key: apiKey,
            max_calls: maxCalls,
            expire_time: expireTime.toISOString(),
            created_at: now.toISOString(),
            provider: provider,
            ext_info: extInfo
        };
    }

    /**
     * Encrypt token data
     * 加密令牌数据
     * @param {object} tokenData Token data / 令牌数据
     * @returns {string} Encrypted token string / 加密后的令牌字符串
     */
    encryptToken(tokenData) {
        // Serialize token data / 序列化令牌数据
        const jsonData = JSON.stringify(tokenData);
        
        // Generate IV / 生成 IV
        const iv = this.generateIV();
        
        // Decode AES key / 解码 AES 密钥
        const key = CryptoJS.enc.Hex.parse(this.config.crypto.aes_key);
        
        // Encrypt data / 加密数据
        const encrypted = CryptoJS.AES.encrypt(jsonData, key, {
            iv: iv,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7
        });
        
        // Combine IV and encrypted data / 组合 IV 和加密数据
        const combined = CryptoJS.lib.WordArray.create()
            .concat(iv)
            .concat(encrypted.ciphertext);
        
        // Base64 URL safe encoding / Base64 URL 安全编码
        return CryptoJS.enc.Base64url.stringify(combined);
    }

    /**
     * Get server URL
     * 获取服务器 URL
     * @param {string} path API path / API 路径
     * @returns {string} Complete server URL / 完整的服务器 URL
     */
    getServerUrl(path = '') {
        const { host, port, base_path } = this.config.server;
        const basePath = base_path.endsWith('/') ? base_path : `${base_path}/`;
        const cleanPath = path.startsWith('/') ? path.slice(1) : path;
        if (path === 'health') {
            return `${host}:${port}/health`;
        }
        return `${host}:${port}${basePath}${cleanPath}`;
    }
} 