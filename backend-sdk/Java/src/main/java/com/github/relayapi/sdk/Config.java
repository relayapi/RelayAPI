package com.github.relayapi.sdk;

import com.fasterxml.jackson.annotation.JsonProperty;

public class Config {
    private String version;
    private ServerConfig server;
    private CryptoConfig crypto;

    public String getVersion() {
        return version;
    }

    public void setVersion(String version) {
        this.version = version;
    }

    public ServerConfig getServer() {
        return server;
    }

    public void setServer(ServerConfig server) {
        this.server = server;
    }

    public CryptoConfig getCrypto() {
        return crypto;
    }

    public void setCrypto(CryptoConfig crypto) {
        this.crypto = crypto;
    }

    public static class ServerConfig {
        private String host;
        private int port;
        @JsonProperty("base_path")
        private String basePath;

        public String getHost() {
            return host;
        }

        public void setHost(String host) {
            this.host = host;
        }

        public int getPort() {
            return port;
        }

        public void setPort(int port) {
            this.port = port;
        }

        public String getBasePath() {
            return basePath;
        }

        public void setBasePath(String basePath) {
            this.basePath = basePath;
        }
    }

    public static class CryptoConfig {
        private String method;
        @JsonProperty("aes_key")
        private String aesKey;
        @JsonProperty("aes_iv_seed")
        private String aesIvSeed;

        public String getMethod() {
            return method;
        }

        public void setMethod(String method) {
            this.method = method;
        }

        public String getAesKey() {
            return aesKey;
        }

        public void setAesKey(String aesKey) {
            this.aesKey = aesKey;
        }

        public String getAesIvSeed() {
            return aesIvSeed;
        }

        public void setAesIvSeed(String aesIvSeed) {
            this.aesIvSeed = aesIvSeed;
        }
    }
} 