import { useState, useRef, useEffect } from 'react'
import type { FormEvent } from 'react'
import { Send } from 'lucide-react'
import { useMutation, useQueryClient } from '@tanstack/react-query'

import { PageShell } from '@/components/layout/PageShell'
import api from '@/lib/api'

interface ChatAction {
    type: string
    name?: string
    amount?: number
    account?: string
    to_account?: string
}

interface ChatMessage {
    role: 'user' | 'assistant'
    content: string
    action?: ChatAction
}

interface ChatApiResponse {
    reply: string
    action?: ChatAction
    success: boolean
}

const SUGGESTIONS = [
    'Spent 500 on groceries',
    'Got 10000 salary',
    'Transfer 2000 from HDFC to SBI',
    'What\'s my balance?',
    'Show recent transactions',
    'Delete the last transaction',
    'help',
]

function actionBadge(action: ChatAction): string {
    switch (action.type) {
        case 'expense': return '🔴 Expense recorded'
        case 'income': return '🟢 Income recorded'
        case 'transfer': return '🔄 Transfer completed'
        case 'delete': return '🗑️ Transaction deleted'
        default: return ''
    }
}

function renderMarkdown(text: string): string {
    return text
        .replace(/\*\*(.+?)\*\*/g, '<strong>$1</strong>')
        .replace(/\n/g, '<br/>')
}

export function ChatPage() {
    const [messages, setMessages] = useState<ChatMessage[]>([
        {
            role: 'assistant',
            content: 'Hi! I\'m your financial assistant. I can help you add transactions, transfer money, check balances, and more.\n\nTry saying things like:\n• "Spent 500 on groceries"\n• "Got 10000 salary"\n• "Transfer 2000 from HDFC to SBI"\n• "What\'s my balance?"\n• "Delete the last transaction"',
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
            if (data.action) {
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
                <div className="flex-1 space-y-4 overflow-y-auto rounded-t-2xl border border-b-0 border-border bg-surface p-4">
                    {messages.map((msg, i) => (
                        <div key={i} className={`flex ${msg.role === 'user' ? 'justify-end' : 'justify-start'}`}>
                            <div
                                className={`max-w-[80%] rounded-2xl px-4 py-3 text-sm ${msg.role === 'user'
                                        ? 'bg-accent text-white'
                                        : 'bg-elevated text-foreground'
                                    }`}
                            >
                                <div dangerouslySetInnerHTML={{ __html: renderMarkdown(msg.content) }} />
                                {msg.action && actionBadge(msg.action) && (
                                    <div className="mt-2 rounded-lg bg-positive-muted px-3 py-1.5 text-xs text-positive">
                                        {actionBadge(msg.action)}
                                    </div>
                                )}
                            </div>
                        </div>
                    ))}
                    {sendMessage.isPending && (
                        <div className="flex justify-start">
                            <div className="rounded-2xl bg-elevated px-4 py-3 text-sm text-muted">
                                Thinking...
                            </div>
                        </div>
                    )}
                    <div ref={messagesEndRef} />
                </div>

                {messages.length <= 1 && (
                    <div className="flex flex-wrap gap-2 border-x border-border bg-surface px-4 py-3">
                        {SUGGESTIONS.map((text) => (
                            <button
                                key={text}
                                type="button"
                                onClick={() => handleSuggestion(text)}
                                className="rounded-full border border-border px-3 py-1.5 text-xs text-secondary hover:bg-surface-hover"
                            >
                                {text}
                            </button>
                        ))}
                    </div>
                )}

                <form
                    onSubmit={handleSubmit}
                    className="flex gap-2 rounded-b-2xl border border-border bg-surface p-3"
                >
                    <input
                        value={input}
                        onChange={(e) => setInput(e.target.value)}
                        placeholder="Type a message... e.g. 'Spent 200 on coffee'"
                        className="flex-1 rounded-lg border border-border bg-elevated px-4 py-2 text-sm text-foreground"
                        disabled={sendMessage.isPending}
                    />
                    <button
                        type="submit"
                        disabled={sendMessage.isPending || !input.trim()}
                        className="rounded-lg bg-accent px-4 py-2 text-white disabled:opacity-60"
                    >
                        <Send size={16} />
                    </button>
                </form>
            </div>
        </PageShell>
    )
}
