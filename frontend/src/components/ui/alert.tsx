import type { HTMLAttributes, ReactNode } from 'react'

type AlertProps = HTMLAttributes<HTMLDivElement> & {
    variant?: 'default' | 'destructive'
    title?: string
    children?: ReactNode
}

export function Alert({ className = '', variant = 'default', title, children, ...props }: AlertProps) {
    const variantClass =
        variant === 'destructive'
            ? 'border-rose-300 bg-rose-50 text-rose-800'
            : 'border-slate-300 bg-slate-50 text-slate-800'

    return (
        <div className={`rounded-lg border p-3 ${variantClass} ${className}`} role="alert" {...props}>
            {title ? <p className="text-sm font-semibold">{title}</p> : null}
            {children ? <div className="text-sm">{children}</div> : null}
        </div>
    )
}
