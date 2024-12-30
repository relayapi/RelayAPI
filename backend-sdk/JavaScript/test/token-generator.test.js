import { TokenGenerator } from '../src/TokenGenerator.js';

describe('TokenGenerator', () => {
    let tokenGenerator;
    const config = {
        version: '1.0.0',
        server: {
            host: 'http://localhost',
            port: 8840,
            base_path: '/relayapi/'
        },
        crypto: {
            method: 'aes',
            aes_key: '0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef',
            aes_iv_seed: 'fedcba9876543210'
        }
    };

    beforeEach(() => {
        tokenGenerator = new TokenGenerator(config);
    });

    describe('initialization', () => {
        it('should initialize with config object', () => {
            expect(tokenGenerator.config).toEqual(config);
        });

        it('should throw error for invalid config', () => {
            expect(() => new TokenGenerator(null)).toThrow('Invalid config parameter');
        });

        it('should throw error for missing crypto section', () => {
            const invalidConfig = { ...config };
            delete invalidConfig.crypto;
            expect(() => new TokenGenerator(invalidConfig)).toThrow('Invalid config: missing crypto or server section');
        });

        it('should throw error for unsupported encryption method', () => {
            const invalidConfig = {
                ...config,
                crypto: { ...config.crypto, method: 'des' }
            };
            expect(() => new TokenGenerator(invalidConfig)).toThrow('Only AES encryption is supported');
        });
    });

    describe('generateIV', () => {
        it('should generate a valid IV', () => {
            const iv = tokenGenerator.generateIV();
            expect(iv).toBeDefined();
            expect(iv.words).toBeDefined();
            expect(iv.words.length).toBeGreaterThan(0);
            expect(iv.sigBytes).toBe(16);
        });
    });

    describe('createToken', () => {
        it('should create token data with default values', () => {
            const tokenData = tokenGenerator.createToken({
                apiKey: 'test-api-key'
            });

            expect(tokenData).toBeDefined();
            expect(tokenData.api_key).toBe('test-api-key');
            expect(tokenData.max_calls).toBe(100);
            expect(tokenData.provider).toBe('dashscope');
            expect(tokenData.ext_info).toBe('');
            expect(new Date(tokenData.expire_time)).toBeInstanceOf(Date);
            expect(new Date(tokenData.created_at)).toBeInstanceOf(Date);
        });

        it('should create token data with custom values', () => {
            const tokenData = tokenGenerator.createToken({
                apiKey: 'test-api-key',
                maxCalls: 200,
                expireDays: 2,
                provider: 'openai',
                extInfo: 'test info'
            });

            expect(tokenData.api_key).toBe('test-api-key');
            expect(tokenData.max_calls).toBe(200);
            expect(tokenData.provider).toBe('openai');
            expect(tokenData.ext_info).toBe('test info');

            const expireTime = new Date(tokenData.expire_time);
            const createdAt = new Date(tokenData.created_at);
            const timeDiff = expireTime.getTime() - createdAt.getTime();
            const daysDiff = Math.round(timeDiff / (1000 * 60 * 60 * 24));
            expect(daysDiff).toBe(2);
        });
    });

    describe('encryptToken', () => {
        it('should encrypt token data', () => {
            const tokenData = tokenGenerator.createToken({
                apiKey: 'test-api-key'
            });

            const encryptedToken = tokenGenerator.encryptToken(tokenData);

            expect(encryptedToken).toBeDefined();
            expect(typeof encryptedToken).toBe('string');
            expect(encryptedToken.length).toBeGreaterThan(0);
            // Base64url 格式验证
            expect(encryptedToken).toMatch(/^[A-Za-z0-9_-]+$/);
        });

        it('should generate different tokens for same data', () => {
            const tokenData = tokenGenerator.createToken({
                apiKey: 'test-api-key'
            });

            const token1 = tokenGenerator.encryptToken(tokenData);
            const token2 = tokenGenerator.encryptToken(tokenData);

            expect(token1).not.toBe(token2);
        });
    });

    describe('getServerUrl', () => {
        it('should generate base URL without path', () => {
            const url = tokenGenerator.getServerUrl();
            expect(url).toBe('http://localhost:8840/relayapi/');
        });

        it('should generate URL with path', () => {
            const url = tokenGenerator.getServerUrl('v1/chat/completions');
            expect(url).toBe('http://localhost:8840/relayapi/v1/chat/completions');
        });

        it('should handle paths with leading slash', () => {
            const url = tokenGenerator.getServerUrl('/v1/chat/completions');
            expect(url).toBe('http://localhost:8840/relayapi/v1/chat/completions');
        });

        it('should handle base paths with and without trailing slash', () => {
            const configWithSlash = {
                ...config,
                server: { ...config.server, base_path: '/relayapi/' }
            };
            const configWithoutSlash = {
                ...config,
                server: { ...config.server, base_path: '/relayapi' }
            };

            const generator1 = new TokenGenerator(configWithSlash);
            const generator2 = new TokenGenerator(configWithoutSlash);

            expect(generator1.getServerUrl('v1/chat')).toBe('http://localhost:8840/relayapi/v1/chat');
            expect(generator2.getServerUrl('v1/chat')).toBe('http://localhost:8840/relayapi/v1/chat');
        });
    });
}); 