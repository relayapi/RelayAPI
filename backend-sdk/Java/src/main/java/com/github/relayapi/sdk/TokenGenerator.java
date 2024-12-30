package com.github.relayapi.sdk;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.bouncycastle.crypto.engines.AESEngine;
import org.bouncycastle.crypto.modes.CBCBlockCipher;
import org.bouncycastle.crypto.paddings.PKCS7Padding;
import org.bouncycastle.crypto.paddings.PaddedBufferedBlockCipher;
import org.bouncycastle.crypto.params.KeyParameter;
import org.bouncycastle.crypto.params.ParametersWithIV;
import org.bouncycastle.util.encoders.Base64;
import org.bouncycastle.util.encoders.Hex;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.SecureRandom;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.HashMap;
import java.util.Map;

public class TokenGenerator {
    private final ObjectMapper objectMapper = new ObjectMapper();
    private final Config config;
    private final String hash;

    public TokenGenerator(Config config) {
        this.config = config;
        validateConfig();
        this.hash = generateConfigHash();
    }

    private void validateConfig() {
        if (config.getCrypto() == null || config.getServer() == null) {
            throw new IllegalArgumentException("Invalid config: missing crypto or server section");
        }
        if (!"aes".equals(config.getCrypto().getMethod())) {
            throw new IllegalArgumentException("Only AES encryption is supported");
        }
        if (config.getCrypto().getAesKey() == null || config.getCrypto().getAesIvSeed() == null) {
            throw new IllegalArgumentException("Invalid crypto config: missing aes_key or aes_iv_seed");
        }
    }

    private String generateConfigHash() {
        try {
            String data = config.getCrypto().getMethod() +
                    config.getCrypto().getAesKey() +
                    config.getCrypto().getAesIvSeed();
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] hash = digest.digest(data.getBytes(StandardCharsets.UTF_8));
            return Hex.toHexString(hash);
        } catch (Exception e) {
            throw new RuntimeException("Failed to generate config hash", e);
        }
    }

    private byte[] generateIV() {
        try {
            byte[] iv = new byte[16];
            SecureRandom random = new SecureRandom();
            random.nextBytes(iv);
            
            // Mix with IV seed
            byte[] ivSeed = config.getCrypto().getAesIvSeed().getBytes(StandardCharsets.UTF_8);
            for (int i = 0; i < iv.length && i < ivSeed.length; i++) {
                iv[i] ^= ivSeed[i];
            }
            return iv;
        } catch (Exception e) {
            throw new RuntimeException("Failed to generate IV", e);
        }
    }

    public Map<String, Object> createToken(TokenOptions options) {
        Instant now = Instant.now();
        Map<String, Object> tokenData = new HashMap<>();
        tokenData.put("id", "token-" + now.toEpochMilli());
        tokenData.put("api_key", options.getApiKey());
        tokenData.put("max_calls", options.getMaxCalls());
        tokenData.put("expire_time", now.plus(options.getExpireSeconds(), ChronoUnit.SECONDS).toString());
        tokenData.put("created_at", now.toString());
        tokenData.put("provider", options.getProvider());
        tokenData.put("ext_info", options.getExtInfo());
        return tokenData;
    }

    public String encryptToken(Map<String, Object> tokenData) {
        try {
            // Convert token data to JSON
            String jsonData = objectMapper.writeValueAsString(tokenData);
            
            // Generate IV
            byte[] iv = generateIV();
            
            // Setup AES cipher
            byte[] key = Hex.decode(config.getCrypto().getAesKey());
            CBCBlockCipher blockCipher = new CBCBlockCipher(new AESEngine());
            PaddedBufferedBlockCipher cipher = new PaddedBufferedBlockCipher(blockCipher, new PKCS7Padding());
            cipher.init(true, new ParametersWithIV(new KeyParameter(key), iv));
            
            // Encrypt data
            byte[] inputBytes = jsonData.getBytes(StandardCharsets.UTF_8);
            byte[] outputBytes = new byte[cipher.getOutputSize(inputBytes.length)];
            int length = cipher.processBytes(inputBytes, 0, inputBytes.length, outputBytes, 0);
            length += cipher.doFinal(outputBytes, length);
            
            // Combine IV and encrypted data
            byte[] combined = new byte[iv.length + length];
            System.arraycopy(iv, 0, combined, 0, iv.length);
            System.arraycopy(outputBytes, 0, combined, iv.length, length);
            
            // Base64URL encode
            return base64UrlEncode(combined);
        } catch (Exception e) {
            throw new RuntimeException("Failed to encrypt token", e);
        }
    }

    private String base64UrlEncode(byte[] data) {
        String base64 = Base64.toBase64String(data);
        return base64.replace('+', '-')
                .replace('/', '_')
                .replace("=", "");
    }

    public String getServerUrl(String path) {
        String basePath = config.getServer().getBasePath();
        if (!basePath.endsWith("/")) {
            basePath += "/";
        }
        String cleanPath = path.startsWith("/") ? path.substring(1) : path;
        
        if ("health".equals(path)) {
            return String.format("%s:%d/health",
                    config.getServer().getHost(),
                    config.getServer().getPort());
        }
        
        return String.format("%s:%d%s%s",
                config.getServer().getHost(),
                config.getServer().getPort(),
                basePath,
                cleanPath);
    }

    public String getHash() {
        return hash;
    }
} 