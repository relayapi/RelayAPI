import CryptoJS from 'crypto-js';
import fs from 'fs/promises';

export class TokenGenerator {
    /**
     * 初始化令牌生成器
     * @param {string|object} config 配置文件路径或配置对象
     */
    constructor(config) {
        this.config = null;
        this.initialize(config);
    }

    /**
     * 初始化配置
     * @param {string|object} config 配置文件路径或配置对象
     */
    async initialize(config) {
        if (typeof config === 'string') {
            try {
                const data = await fs.readFile(config, 'utf8');
                this.config = JSON.parse(data);
            } catch (error) {
                if (config === 'default.rai') {
                    // 使用默认配置
                    this.config = {
                        version: '1.0.0',
                        server: {
                            host: 'http://localhost',
                            port: 8080,
                            base_path: '/relayapi/'
                        },
                        crypto: {
                            method: 'aes',
                            aes_key: '0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef',
                            aes_iv_seed: 'fedcba9876543210'
                        }
                    };
                } else {
                    throw new Error(`Failed to load config: ${error.message}`);
                }
            }
        } else if (typeof config === 'object') {
            this.config = config;
        } else {
            throw new Error('Invalid config parameter');
        }

        // 验证配置
        this.validateConfig();
    }

    /**
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
     * 生成初始化向量
     * @returns {WordArray} 生成的 IV
     */
    generateIV() {
        // 生成随机 IV
        const iv = CryptoJS.lib.WordArray.random(16);
        // 使用 IV 种子进行混合
        const ivSeed = CryptoJS.enc.Utf8.parse(this.config.crypto.aes_iv_seed);
        const words = iv.words.map((word, i) => word ^ ivSeed.words[i]);
        return CryptoJS.lib.WordArray.create(words, 16);
    }

    /**
     * 创建令牌数据
     * @param {object} options 令牌选项
     * @returns {object} 令牌数据
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
     * 加密令牌数据
     * @param {object} tokenData 令牌数据
     * @returns {string} 加密后的令牌字符串
     */
    encryptToken(tokenData) {
        // 序列化令牌数据
        const jsonData = JSON.stringify(tokenData);
        
        // 生成 IV
        const iv = this.generateIV();
        
        // 解码 AES 密钥
        const key = CryptoJS.enc.Hex.parse(this.config.crypto.aes_key);
        
        // 加密数据
        const encrypted = CryptoJS.AES.encrypt(jsonData, key, {
            iv: iv,
            mode: CryptoJS.mode.CBC,
            padding: CryptoJS.pad.Pkcs7
        });
        
        // 组合 IV 和加密数据
        const combined = CryptoJS.lib.WordArray.create()
            .concat(iv)
            .concat(encrypted.ciphertext);
        
        // Base64 URL 安全编码
        return CryptoJS.enc.Base64url.stringify(combined);
    }

    /**
     * 获取服务器 URL
     * @param {string} path API 路径
     * @returns {string} 完整的服务器 URL
     */
    getServerUrl(path = '') {
        const { host, port, base_path } = this.config.server;
        const basePath = base_path.endsWith('/') ? base_path : `${base_path}/`;
        const cleanPath = path.startsWith('/') ? path.slice(1) : path;
        return `${host}:${port}${basePath}${cleanPath}`;
    }
} 