import * as AccordionPrimitive from '@radix-ui/react-accordion'
import { ChevronDown } from 'lucide-react'
import { cn } from '@/lib/utils'

const Accordion = AccordionPrimitive.Root

const AccordionItem = ({
    className,
    ...props
}: React.ComponentPropsWithoutRef<typeof AccordionPrimitive.Item>) => (
    <AccordionPrimitive.Item className={cn('border-b border-slate-100', className)} {...props} />
)
AccordionItem.displayName = 'AccordionItem'

const AccordionTrigger = ({
    className,
    children,
    ...props
}: React.ComponentPropsWithoutRef<typeof AccordionPrimitive.Trigger>) => (
    <AccordionPrimitive.Header className="flex">
        <AccordionPrimitive.Trigger
            className={cn(
                'flex flex-1 items-center justify-between py-4 text-sm font-semibold text-slate-900 transition-all hover:underline [&[data-state=open]>svg]:rotate-180',
                className,
            )}
            {...props}
        >
            {children}
            <ChevronDown className="h-4 w-4 shrink-0 text-slate-500 transition-transform duration-200" />
        </AccordionPrimitive.Trigger>
    </AccordionPrimitive.Header>
)
AccordionTrigger.displayName = AccordionPrimitive.Trigger.displayName

const AccordionContent = ({
    className,
    children,
    ...props
}: React.ComponentPropsWithoutRef<typeof AccordionPrimitive.Content>) => (
    <AccordionPrimitive.Content
        className={cn(
            'overflow-hidden text-sm data-[state=closed]:animate-accordion-up data-[state=open]:animate-accordion-down',
            className,
        )}
        {...props}
    >
        <div className={cn('pb-4 pt-0 leading-6 text-slate-600')}>{children}</div>
    </AccordionPrimitive.Content>
)
AccordionContent.displayName = AccordionPrimitive.Content.displayName

export { Accordion, AccordionItem, AccordionTrigger, AccordionContent }
