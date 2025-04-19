<template>
    <div class="chat-container">
      <div class="chat-header" @click="toggleChat">
        {{ isOpen ? "✖ CLOSE" : "💬 Chat with Gemini" }}
      </div>
      <div v-if="isOpen" class="chat-box">
        <div class="messages">
          <div v-for="(message, index) in messages" :key="index" :class="['message', message.role]">
            <span>{{ message.text }}</span>
          </div>
        </div>
        <input
          v-model="userInput"
          @keyup.enter="sendMessage"
          placeholder="Type your question..."
          class="chat-input"
          :disabled="loading"
        />
        <button v-if="loading" class="loading-btn">sending...</button>
      </div>
    </div>
  </template>
  
  <script>
  import axios from "axios";
  
  export default {
    data() {
      return {
        isOpen: false,
        userInput: "",
        messages: [],
        loading: false,
        apiUrl: "http://170.106.117.71:3000/chat" // ✅ 直接请求 Node.js 后端
      };
    },
    methods: {
      toggleChat() {
        this.isOpen = !this.isOpen;
      },
  
      async sendMessage() {
        if (!this.userInput.trim()) return;

        this.messages.push({ role: "user", text: this.userInput });

        this.loading = true;
        const userInputText = this.userInput;
        this.userInput = "";

        try {
          console.log("📡 发送请求到:", this.apiUrl);

          const response = await axios.post(this.apiUrl, { message: userInputText });

          console.log("🌍 API 完整响应:", response); // ✅ 确保 Vue 收到了完整的 HTTP 响应
          console.log("🌍 API 响应数据:", response.reply); // ✅ 确保 `response.data` 不为空

          if (!response || !response.reply) {
            throw new Error("API 响应为空");
          }

          const botResponse =
            response.reply || "未收到 Gemini 响应";

          this.messages.push({ role: "bot", text: botResponse });

        } catch (error) {
          console.error("❌ Gemini API 错误:", error);
          this.messages.push({ role: "bot", text: "请求失败，请稍后再试" });

        } finally {
          this.loading = false;
        }
      }


    }
  };
  </script>
  
  <style scoped>
  .chat-container {
    position: fixed;
    bottom: 50px;
    left: 20px;
    width: 300px;
    font-family: Arial, sans-serif;
  }
  
  .chat-header {
    background-color: #000000;
    color: white;
    padding: 10px;
    cursor: pointer;
    text-align: center; /* 让文字靠右 */
    border-radius: 10px 10px 10px 10px;
}
  
  .chat-box {
    background: white;
    border: 1px solid #ddd;
    border-top: none;
    border-radius: 10px 10px 10px 10px;
    overflow: hidden;
    display: flex;
    flex-direction: column;
  }
  
  .messages {
    max-height: 300px;
    overflow-y: auto;
    padding: 10px;
    display: flex;
    flex-direction: column;
  }
  
  .message {
    padding: 8px;
    border-radius: 5px;
    margin-bottom: 5px;
    max-width: 80%;
  }
  
  .message.user {
    background-color: #000000; /* 黑色背景 */
    color: white; /* 文字变成白色 */
    align-self: flex-end;
}
  
  .message.bot {
    background-color: #f1f1f1;
    color: black;
    align-self: flex-start;
  }
  
  .chat-input {
    border: none;
    padding: 10px;
    width: calc(100% - 20px);
    margin: 10px;
    border-radius: 5px;
    outline: none;
  }
  
  .loading-btn {
    background-color: gray;
    color: white;
    padding: 8px;
    border: none;
    cursor: not-allowed;
    text-align: center;
  }
  </style>
  
