import type { ButtonHTMLAttributes } from 'react'

type ButtonVariant = 'default' | 'outline'

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: ButtonVariant
}

export function Button({ className = '', variant = 'default', ...props }: ButtonProps) {
    const variantClass =
        variant === 'outline'
            ? 'border border-border bg-surface text-foreground hover:bg-surface-hover'
            : 'bg-accent text-white hover:bg-accent/80'

    return (
        <button
            className={`inline-flex items-center justify-center rounded-lg px-4 py-2 text-sm font-medium transition disabled:cursor-not-allowed disabled:opacity-60 ${variantClass} ${className}`}
            {...props}
        />
    )
}
