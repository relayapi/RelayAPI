package com.github.relayapi.sdk;

import com.fasterxml.jackson.databind.ObjectMapper;
import okhttp3.*;

import java.io.IOException;
import java.util.List;
import java.util.Map;
import java.util.concurrent.TimeUnit;

public class RelayAPIClient {
    private static final MediaType JSON = MediaType.parse("application/json; charset=utf-8");
    private final TokenGenerator tokenGenerator;
    private final OkHttpClient httpClient;
    private final ObjectMapper objectMapper;

    public RelayAPIClient(Config config) {
        this.tokenGenerator = new TokenGenerator(config);
        this.httpClient = new OkHttpClient.Builder()
                .connectTimeout(30, TimeUnit.SECONDS)
                .readTimeout(30, TimeUnit.SECONDS)
                .writeTimeout(30, TimeUnit.SECONDS)
                .build();
        this.objectMapper = new ObjectMapper();
    }

    public String createToken(TokenOptions options) {
        Map<String, Object> tokenData = tokenGenerator.createToken(options);
        return tokenGenerator.encryptToken(tokenData);
    }

    public String generateUrl(String token, String endpoint) {
        String baseUrl = tokenGenerator.getServerUrl(endpoint);
        return String.format("%s?token=%s&rai_hash=%s", baseUrl, token, tokenGenerator.getHash());
    }

    public Map<String, Object> chat(ChatRequest request) throws IOException {
        String url = generateUrl(request.getToken(), "chat/completions");
        RequestBody body = RequestBody.create(
                objectMapper.writeValueAsString(request.toRequestBody()),
                JSON
        );

        Request httpRequest = new Request.Builder()
                .url(url)
                .post(body)
                .build();

        try (Response response = httpClient.newCall(httpRequest).execute()) {
            if (!response.isSuccessful()) {
                throw new IOException("Chat request failed: " + response.code());
            }
            return objectMapper.readValue(response.body().string(), Map.class);
        }
    }

    public Map<String, Object> generateImage(ImageRequest request) throws IOException {
        String url = generateUrl(request.getToken(), "images/generations");
        RequestBody body = RequestBody.create(
                objectMapper.writeValueAsString(request.toRequestBody()),
                JSON
        );

        Request httpRequest = new Request.Builder()
                .url(url)
                .post(body)
                .build();

        try (Response response = httpClient.newCall(httpRequest).execute()) {
            if (!response.isSuccessful()) {
                throw new IOException("Image generation failed: " + response.code());
            }
            return objectMapper.readValue(response.body().string(), Map.class);
        }
    }

    public Map<String, Object> createEmbedding(EmbeddingRequest request) throws IOException {
        String url = generateUrl(request.getToken(), "embeddings");
        RequestBody body = RequestBody.create(
                objectMapper.writeValueAsString(request.toRequestBody()),
                JSON
        );

        Request httpRequest = new Request.Builder()
                .url(url)
                .post(body)
                .build();

        try (Response response = httpClient.newCall(httpRequest).execute()) {
            if (!response.isSuccessful()) {
                throw new IOException("Embedding creation failed: " + response.code());
            }
            return objectMapper.readValue(response.body().string(), Map.class);
        }
    }

    public Map<String, Object> healthCheck() throws IOException {
        String url = tokenGenerator.getServerUrl("health");
        Request request = new Request.Builder()
                .url(url)
                .get()
                .build();

        try (Response response = httpClient.newCall(request).execute()) {
            if (!response.isSuccessful()) {
                throw new IOException("Health check failed: " + response.code());
            }
            return objectMapper.readValue(response.body().string(), Map.class);
        }
    }

    public static class ChatRequest {
        private List<Map<String, String>> messages;
        private String model = "gpt-3.5-turbo";
        private double temperature = 0.7;
        private int maxTokens = 1000;
        private String token;

        public ChatRequest(String token) {
            this.token = token;
        }

        public List<Map<String, String>> getMessages() {
            return messages;
        }

        public ChatRequest setMessages(List<Map<String, String>> messages) {
            this.messages = messages;
            return this;
        }

        public String getModel() {
            return model;
        }

        public ChatRequest setModel(String model) {
            this.model = model;
            return this;
        }

        public double getTemperature() {
            return temperature;
        }

        public ChatRequest setTemperature(double temperature) {
            this.temperature = temperature;
            return this;
        }

        public int getMaxTokens() {
            return maxTokens;
        }

        public ChatRequest setMaxTokens(int maxTokens) {
            this.maxTokens = maxTokens;
            return this;
        }

        public String getToken() {
            return token;
        }

        public Map<String, Object> toRequestBody() {
            return Map.of(
                    "messages", messages,
                    "model", model,
                    "temperature", temperature,
                    "max_tokens", maxTokens
            );
        }
    }

    public static class ImageRequest {
        private String prompt;
        private String model = "dall-e-3";
        private String size = "1024x1024";
        private String quality = "standard";
        private int n = 1;
        private String token;

        public ImageRequest(String token) {
            this.token = token;
        }

        public String getPrompt() {
            return prompt;
        }

        public ImageRequest setPrompt(String prompt) {
            this.prompt = prompt;
            return this;
        }

        public String getModel() {
            return model;
        }

        public ImageRequest setModel(String model) {
            this.model = model;
            return this;
        }

        public String getSize() {
            return size;
        }

        public ImageRequest setSize(String size) {
            this.size = size;
            return this;
        }

        public String getQuality() {
            return quality;
        }

        public ImageRequest setQuality(String quality) {
            this.quality = quality;
            return this;
        }

        public int getN() {
            return n;
        }

        public ImageRequest setN(int n) {
            this.n = n;
            return this;
        }

        public String getToken() {
            return token;
        }

        public Map<String, Object> toRequestBody() {
            return Map.of(
                    "prompt", prompt,
                    "model", model,
                    "size", size,
                    "quality", quality,
                    "n", n
            );
        }
    }

    public static class EmbeddingRequest {
        private String input;
        private String model = "text-embedding-ada-002";
        private String token;

        public EmbeddingRequest(String token) {
            this.token = token;
        }

        public String getInput() {
            return input;
        }

        public EmbeddingRequest setInput(String input) {
            this.input = input;
            return this;
        }

        public String getModel() {
            return model;
        }

        public EmbeddingRequest setModel(String model) {
            this.model = model;
            return this;
        }

        public String getToken() {
            return token;
        }

        public Map<String, Object> toRequestBody() {
            return Map.of(
                    "input", input,
                    "model", model
            );
        }
    }
} 