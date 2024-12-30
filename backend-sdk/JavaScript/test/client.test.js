import { RelayAPIClient } from '../src/Client.js';

describe('RelayAPIClient', () => {
    let client;
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
        client = new RelayAPIClient(config);
    });

    describe('createToken', () => {
        it('should create a valid token', () => {
            const token = client.createToken({
                apiKey: 'test-api-key',
                maxCalls: 100,
                expireDays: 1,
                provider: 'openai'
            });

            expect(token).toBeDefined();
            expect(typeof token).toBe('string');
            expect(token.length).toBeGreaterThan(0);
        });
    });

    describe('generateUrl', () => {
        it('should generate a valid URL with token', () => {
            const token = client.createToken({
                apiKey: 'test-api-key'
            });

            const url = client.generateUrl('v1/chat/completions', token);

            expect(url).toBeDefined();
            expect(url).toContain('http://localhost:8840/relayapi/v1/chat/completions');
            expect(url).toContain('token=');
        });
    });

    describe('chat', () => {
        it('should make a chat request', async () => {
            const token = client.createToken({
                apiKey: 'test-api-key'
            });

            // Mock axios post request
            const mockResponse = {
                data: {
                    choices: [
                        {
                            message: {
                                content: 'Hello, how can I help you?'
                            }
                        }
                    ]
                }
            };

            global.fetch = jest.fn(() =>
                Promise.resolve({
                    ok: true,
                    json: () => Promise.resolve(mockResponse.data)
                })
            );

            const response = await client.chat({
                messages: [
                    { role: 'user', content: 'Hello!' }
                ],
                token
            });

            expect(response).toBeDefined();
            expect(response.choices[0].message.content).toBe('Hello, how can I help you?');
        });
    });

    describe('healthCheck', () => {
        it('should check server health status', async () => {
            const mockResponse = {
                data: {
                    status: 'ok',
                    version: '1.0.0'
                }
            };

            global.fetch = jest.fn(() =>
                Promise.resolve({
                    ok: true,
                    json: () => Promise.resolve(mockResponse.data)
                })
            );

            const status = await client.healthCheck();

            expect(status).toBeDefined();
            expect(status.status).toBe('ok');
            expect(status.version).toBe('1.0.0');
        });
    });
}); 