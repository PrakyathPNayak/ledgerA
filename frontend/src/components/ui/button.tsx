import type { ButtonHTMLAttributes } from 'react'

type ButtonVariant = 'default' | 'outline'

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: ButtonVariant
}

export function Button({ className = '', variant = 'default', ...props }: ButtonProps) {
    const variantClass =
        variant === 'outline'
            ? 'border border-slate-300 bg-white text-slate-900 hover:bg-slate-100'
            : 'bg-slate-900 text-white hover:bg-slate-700'

    return (
        <button
            className={`inline-flex items-center justify-center rounded-lg px-4 py-2 text-sm font-medium transition disabled:cursor-not-allowed disabled:opacity-60 ${variantClass} ${className}`}
            {...props}
        />
    )
}
