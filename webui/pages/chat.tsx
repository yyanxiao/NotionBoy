import { useEffect, useRef, useState } from "react"
import Head from "next/head"
import { marked } from "marked"

import { CreateConversationRequest } from "@/lib/pb/model/conversation.pb"
import { Service } from "@/lib/pb/server.pb"

function parseMarkdown(text: string) {
  return { __html: marked(text) }
}

export default function ChatPage() {
  const [messages, setMessages] = useState([])

  const handleSubmit = (e) => {
    if (input === "") return
    e.preventDefault()
    setInput("")
    const req: CreateConversationRequest = {}
    Service.CreateConversation(req).then((res) => {
      console.log(res)
    })

    Service.Status("").then((res) => {
      console.log(res)
    })

    setMessages((messages) => [
      ...messages,
      { id: messages.length, text: input },
    ])
    setInput("")
  }

  const [input, setInput] = useState("")
  const messagesRef = useRef(null)
  useEffect(() => {
    messagesRef.current.scrollTop = messagesRef.current.scrollHeight
  }, [messages])

  return (
    <div className="flex flex-col h-screen">
      <Head>
        <title>Chat</title>
      </Head>
      <main
        ref={messagesRef}
        className="flex-1 overflow-auto w-full flex flex-col"
      >
        <div className="flex-1 border border-blue-400 m-4">
          <div className="flex flex-col">
            {messages.map((message) => (
              <div
                key={message.id}
                className={`prose text-sm max-w-none p-1 rounded-lg shadow-lg mb-4 word-warp break-words
                ${
                  message.id % 2 === 0
                    ? "mr-5  bg-green-100 md:mr-10"
                    : "ml-5  bg-blue-100 md:ml-10 "
                }`}
                dangerouslySetInnerHTML={parseMarkdown(message.text)}
              ></div>
            ))}
          </div>
        </div>
      </main>
      <div
        className="
          w-full
          h-24
          border border-blue-400
          flex flex-row
          "
      >
        <textarea
          placeholder="Type your message here. Press Shift + Enter to send."
          className="stretch flex-1 h-24 border border-blue-400 p-2"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && e.shiftKey) {
              handleSubmit(e)
            }
          }}
        ></textarea>
        <button
          className="w-20 h-24 border bg-red-200"
          onClick={(e) => handleSubmit(e)}
        >
          Send Message
        </button>
      </div>
    </div>
  )
}
