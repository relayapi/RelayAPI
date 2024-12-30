package com.github.relayapi.sdk.examples;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.github.relayapi.sdk.Config;
import com.github.relayapi.sdk.RelayAPIClient;
import com.github.relayapi.sdk.TokenOptions;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.List;
import java.util.Map;
import java.util.Properties;

public class ChatExample {
    public static void main(String[] args) {
        try {
            // 从 .env 文件加载 API key
            Properties env = new Properties();
            env.load(Files.newBufferedReader(Paths.get(".env")));
            String apiKey = env.getProperty("RELAY_API_KEY");
            if (apiKey == null || apiKey.trim().isEmpty()) {
                throw new IllegalStateException("RELAY_API_KEY not found in .env file");
            }

            // 从文件加载配置
            String configContent = Files.readString(Paths.get("config.rai"));
            ObjectMapper mapper = new ObjectMapper();
            Config config = mapper.readValue(configContent, Config.class);

            // 创建客户端实例
            RelayAPIClient client = new RelayAPIClient(config);

            // 创建令牌
            TokenOptions options = new TokenOptions(apiKey)
                    .setMaxCalls(100)
                    .setExpireSeconds(3600)
                    .setProvider("dashscope");
            String token = client.createToken(options);

            // 发送聊天请求
            RelayAPIClient.ChatRequest chatRequest = new RelayAPIClient.ChatRequest(token)
                    .setMessages(List.of(
                            Map.of("role", "system", "content", "You are a helpful assistant."),
                            Map.of("role", "user", "content", "What is the capital of France?")
                    ))
                    .setModel("qwen-vl-max")
                    .setTemperature(0.7)
                    .setMaxTokens(1000);

            System.out.println("Sending chat request...");
            Map<String, Object> chatResponse = client.chat(chatRequest);
            System.out.println("Chat response: " + mapper.writerWithDefaultPrettyPrinter().writeValueAsString(chatResponse));

            // 生成图像
            RelayAPIClient.ImageRequest imageRequest = new RelayAPIClient.ImageRequest(token)
                    .setPrompt("A beautiful sunset over Paris")
                    .setModel("dall-e-3")
                    .setSize("1024x1024")
                    .setQuality("standard")
                    .setN(1);

            System.out.println("\nGenerating image...");
            Map<String, Object> imageResponse = client.generateImage(imageRequest);
            System.out.println("Image response: " + mapper.writerWithDefaultPrettyPrinter().writeValueAsString(imageResponse));

            // 生成嵌入向量
            RelayAPIClient.EmbeddingRequest embeddingRequest = new RelayAPIClient.EmbeddingRequest(token)
                    .setInput("The quick brown fox jumps over the lazy dog")
                    .setModel("text-embedding-ada-002");

            System.out.println("\nGenerating embedding...");
            Map<String, Object> embeddingResponse = client.createEmbedding(embeddingRequest);
            System.out.println("Embedding response: " + mapper.writerWithDefaultPrettyPrinter().writeValueAsString(embeddingResponse));

            // 健康检查
            System.out.println("\nChecking server health...");
            Map<String, Object> healthStatus = client.healthCheck();
            System.out.println("Health status: " + mapper.writerWithDefaultPrettyPrinter().writeValueAsString(healthStatus));

        } catch (Exception e) {
            System.err.println("Error: " + e.getMessage());
            e.printStackTrace();
        }
    }
} 