

require("dotenv").config();
const express = require("express");
const axios = require("axios");

const cors = require("cors");



const app = express();
app.use(cors({
    origin: "*",  // 或者 "http://localhost:8000"
    methods: ["POST"],
    allowedHeaders: ["Content-Type"]
  }));
const PORT = 3000;

app.use(express.json());

const GEMINI_API_KEY = "AIzaSyASFi2GlRG72jYFCqoXltB6z0qtSlaQnnA";
const GEMINI_API_URL = `https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=${GEMINI_API_KEY}`;

app.post("/chat", async (req, res) => {
    const userMessage = req.body.message;
    if (!userMessage) {
        return res.status(400).json({ error: "Message is required" });
    }

    try {
        const response = await axios.post(GEMINI_API_URL, {
            contents: [{ parts: [{ text: userMessage }] }]
        });

        const reply = response.data?.candidates?.[0]?.content?.parts?.[0]?.text || "Gemini 没有返回结果";
        res.json({ reply });
    } catch (error) {
        console.error("Gemini API Error:", error.response?.data || error.message);
        res.status(500).json({ error: "服务器内部错误，请检查后端日志" });
    }
});

app.listen(PORT, () => {
    console.log(`Server running at http://localhost:${PORT}`);
});