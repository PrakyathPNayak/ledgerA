import { useState, useRef, useEffect } from 'react'
import type { FormEvent } from 'react'
import { Send } from 'lucide-react'
import { useMutation, useQueryClient } from '@tanstack/react-query'

import { PageShell } from '@/components/layout/PageShell'
import api from '@/lib/api'

interface ChatMessage {
    role: 'user' | 'assistant'
    content: string
    action?: { type: string; name?: string; amount?: number }
}

interface ChatApiResponse {
    reply: string
    action?: { type: string; name?: string; amount?: number }
    success: boolean
}

const SUGGESTIONS = [
    'Spent 500 on groceries',
    'Got 10000 salary',
    'What\'s my balance?',
    'Show recent transactions',
    'help',
]

export function ChatPage() {
    const [messages, setMessages] = useState<ChatMessage[]>([
        {
            role: 'assistant',
            content: 'Hi! I\'m your financial assistant. I can help you add transactions, check balances, and more.\n\nTry saying things like:\n• "Spent 500 on groceries"\n• "Got 10000 salary"\n• "What\'s my balance?"\n• "Show recent transactions"',
        },
    ])
    const [input, setInput] = useState('')
    const messagesEndRef = useRef<HTMLDivElement>(null)
    const queryClient = useQueryClient()

    const sendMessage = useMutation({
        mutationFn: (message: string) =>
            api.post('/chat', { message }) as Promise<ChatApiResponse>,
        onSuccess: (data) => {
            setMessages((prev) => [...prev, { role: 'assistant', content: data.reply, action: data.action }])
            if (data.action?.type === 'expense' || data.action?.type === 'income') {
                queryClient.invalidateQueries({ queryKey: ['transactions'] })
                queryClient.invalidateQueries({ queryKey: ['accounts'] })
                queryClient.invalidateQueries({ queryKey: ['stats'] })
            }
        },
        onError: () => {
            setMessages((prev) => [
                ...prev,
                { role: 'assistant', content: 'Sorry, something went wrong. Please try again.' },
            ])
        },
    })

    useEffect(() => {
        messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
    }, [messages])

    function handleSubmit(e: FormEvent) {
        e.preventDefault()
        const trimmed = input.trim()
        if (!trimmed || sendMessage.isPending) return

        setMessages((prev) => [...prev, { role: 'user', content: trimmed }])
        setInput('')
        sendMessage.mutate(trimmed)
    }

    function handleSuggestion(text: string) {
        if (sendMessage.isPending) return
        setMessages((prev) => [...prev, { role: 'user', content: text }])
        sendMessage.mutate(text)
    }

    return (
        <PageShell title="Chat">
            <div className="mx-auto flex h-[calc(100vh-10rem)] max-w-3xl flex-col">
                <div className="flex-1 space-y-4 overflow-y-auto rounded-t-2xl border border-b-0 border-slate-200 bg-white p-4 dark:border-slate-700 dark:bg-slate-900">
                    {messages.map((msg, i) => (
                        <div key={i} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                            <div
                                className={`max-w-[80%] rounded-2xl px-4 py-3 text-sm whitespace-pre-wrap ${msg.role === 'user'
                                        ? 'bg-slate-900 text-white dark:bg-slate-100 dark:text-slate-900'
                                        : 'bg-slate-100 text-slate-900 dark:bg-slate-800 dark:text-slate-100'
                                    }`}
                            >
                                {msg.content}
                                {msg.action && (
                                    <div className="mt-2 rounded-lg bg-emerald-100 px-3 py-1.5 text-xs text-emerald-800 dark:bg-emerald-900 dark:text-emerald-200">
                                        {msg.action.type === 'expense' ? '🔴 Expense' : '🟢 Income'} recorded
                                    </div>
                                )}
                            </div>
                        </div>
                    ))}
                    {sendMessage.isPending && (
                        <div className="flex justify-start">
                            <div className="rounded-2xl bg-slate-100 px-4 py-3 text-sm text-slate-500 dark:bg-slate-800 dark:text-slate-400">
                                Thinking...
                            </div>
                        </div>
                    )}
                    <div ref={messagesEndRef} />
                </div>

                {messages.length <= 1 && (
                    <div className="flex flex-wrap gap-2 border-x border-slate-200 bg-white px-4 py-3 dark:border-slate-700 dark:bg-slate-900">
                        {SUGGESTIONS.map((text) => (
                            <button
                                key={text}
                                type="button"
                                onClick={() => handleSuggestion(text)}
                                className="rounded-full border border-slate-300 px-3 py-1.5 text-xs text-slate-700 hover:bg-slate-50 dark:border-slate-600 dark:text-slate-300 dark:hover:bg-slate-800"
                            >
                                {text}
                            </button>
                        ))}
                    </div>
                )}

                <form
                    onSubmit={handleSubmit}
                    className="flex gap-2 rounded-b-2xl border border-slate-200 bg-white p-3 dark:border-slate-700 dark:bg-slate-900"
                >
                    <input
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        placeholder="Type a message... e.g. 'Spent 200 on coffee'"
                        className="flex-1 rounded-lg border border-slate-300 px-4 py-2 text-sm dark:border-slate-700 dark:bg-slate-800 dark:text-slate-100"
                        disabled={sendMessage.isPending}
                    />
                    <button
                        type="submit"
                        disabled={sendMessage.isPending || !input.trim()}
                        className="rounded-lg bg-slate-900 px-4 py-2 text-white disabled:opacity-60 dark:bg-slate-100 dark:text-slate-900"
                    >
                        <Send size={16} />
                    </button>
                </form>
            </div>
        </PageShell>
    )
}
