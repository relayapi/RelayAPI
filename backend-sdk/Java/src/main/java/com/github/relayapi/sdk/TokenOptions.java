package com.github.relayapi.sdk;

public class TokenOptions {
    private String apiKey;
    private int maxCalls = 100;
    private long expireSeconds = 86400;
    private String provider = "dashscope";
    private String extInfo = "";

    public TokenOptions(String apiKey) {
        this.apiKey = apiKey;
    }

    public String getApiKey() {
        return apiKey;
    }

    public void setApiKey(String apiKey) {
        this.apiKey = apiKey;
    }

    public int getMaxCalls() {
        return maxCalls;
    }

    public TokenOptions setMaxCalls(int maxCalls) {
        this.maxCalls = maxCalls;
        return this;
    }

    public long getExpireSeconds() {
        return expireSeconds;
    }

    public TokenOptions setExpireSeconds(long expireSeconds) {
        this.expireSeconds = expireSeconds;
        return this;
    }

    public String getProvider() {
        return provider;
    }

    public TokenOptions setProvider(String provider) {
        this.provider = provider;
        return this;
    }

    public String getExtInfo() {
        return extInfo;
    }

    public TokenOptions setExtInfo(String extInfo) {
        this.extInfo = extInfo;
        return this;
    }
} 