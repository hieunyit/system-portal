
"use client"

import * as React from "react"
import * as ToastPrimitives from "@radix-ui/react-toast"
import { cva, type VariantProps } from "class-variance-authority"
import { X, CheckCircle, AlertTriangle, Info, Bell } from "lucide-react"

import { cn } from "@/lib/utils"

const ToastProvider = ToastPrimitives.Provider

const ToastViewport = React.forwardRef<
  React.ElementRef<typeof ToastPrimitives.Viewport>,
  React.ComponentPropsWithoutRef<typeof ToastPrimitives.Viewport>
>(({ className, ...props }, ref) => (
  <ToastPrimitives.Viewport
    ref={ref}
    className={cn(
      "fixed top-0 z-[100] flex max-h-screen w-full flex-col-reverse p-4 sm:bottom-0 sm:right-0 sm:top-auto sm:flex-col md:max-w-[420px]",
      className
    )}
    {...props}
  />
))
ToastViewport.displayName = ToastPrimitives.Viewport.displayName

const toastVariants = cva(
  "group pointer-events-auto relative flex w-full items-start space-x-3 overflow-hidden rounded-md border p-4 pr-8 shadow-lg transition-all data-[swipe=cancel]:translate-x-0 data-[swipe=end]:translate-x-[var(--radix-toast-swipe-end-x)] data-[swipe=move]:translate-x-[var(--radix-toast-swipe-move-x)] data-[swipe=move]:transition-none data-[state=closed]:fade-out-80 data-[state=closed]:slide-out-to-right-full data-[state=open]:sm:slide-in-from-bottom-full",
  // Removed animation classes: data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=open]:slide-in-from-top-full
  {
    variants: {
      variant: {
        default: "border bg-background text-foreground [&_svg]:text-foreground",
        destructive:
          "destructive group border-destructive bg-destructive text-destructive-foreground dark:border-destructive [&_svg]:text-destructive-foreground",
        success:
          "success group border-green-500/50 bg-green-500/10 text-green-700 dark:border-green-600 dark:bg-green-500/20 dark:text-green-400 [&_svg]:text-green-600 dark:[&_svg]:text-green-500",
        info:
          "info group border-blue-500/50 bg-blue-500/10 text-blue-700 dark:border-blue-600 dark:bg-blue-500/20 dark:text-blue-400 [&_svg]:text-blue-600 dark:[&_svg]:text-blue-500",
        warning:
          "warning group border-yellow-500/50 bg-yellow-500/10 text-yellow-700 dark:border-yellow-600 dark:bg-yellow-500/20 dark:text-yellow-400 [&_svg]:text-yellow-600 dark:[&_svg]:text-yellow-500",
      },
    },
    defaultVariants: {
      variant: "default",
    },
  }
)

const Toast = React.forwardRef<
  React.ElementRef<typeof ToastPrimitives.Root>,
  React.ComponentPropsWithoutRef<typeof ToastPrimitives.Root> &
    VariantProps<typeof toastVariants> & {
      icon?: React.ReactNode
    }
>(({ className, variant, children, icon, ...props }, ref) => {
  const defaultIcon =
    variant === "success" ? <CheckCircle className="h-5 w-5" /> :
    variant === "destructive" ? <AlertTriangle className="h-5 w-5" /> :
    variant === "info" ? <Info className="h-5 w-5" /> :
    variant === "warning" ? <Bell className="h-5 w-5" /> :
    null
  
  const displayIcon = icon !== undefined ? icon : defaultIcon;

  return (
    <ToastPrimitives.Root
      ref={ref}
      className={cn(toastVariants({ variant }), className)}
      {...props}
    >
      {displayIcon && <div className="shrink-0 pt-0.5">{displayIcon}</div>}
      <div className="flex-1">{children}</div>
    </ToastPrimitives.Root>
  )
})
Toast.displayName = ToastPrimitives.Root.displayName

const ToastAction = React.forwardRef<
  React.ElementRef<typeof ToastPrimitives.Action>,
  React.ComponentPropsWithoutRef<typeof ToastPrimitives.Action>
>(({ className, ...props }, ref) => (
  <ToastPrimitives.Action
    ref={ref}
    className={cn(
      "inline-flex h-8 shrink-0 items-center justify-center rounded-md border bg-transparent px-3 text-sm font-medium ring-offset-background transition-colors hover:bg-secondary focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
      "group-[.destructive]:border-muted/40 group-[.destructive]:hover:border-destructive/30 group-[.destructive]:hover:bg-destructive-foreground group-[.destructive]:hover:text-destructive group-[.destructive]:focus:ring-destructive", 
      "group-[.success]:border-muted/40 group-[.success]:hover:border-green-500/30 group-[.success]:hover:bg-green-600 group-[.success]:hover:text-primary-foreground group-[.success]:focus:ring-green-600",
      "group-[.info]:border-muted/40 group-[.info]:hover:border-blue-500/30 group-[.info]:hover:bg-blue-600 group-[.info]:hover:text-primary-foreground group-[.info]:focus:ring-blue-600",
      "group-[.warning]:border-muted/40 group-[.warning]:hover:border-yellow-500/30 group-[.warning]:hover:bg-yellow-600 group-[.warning]:hover:text-primary-foreground group-[.warning]:focus:ring-yellow-600",
      className
    )}
    {...props}
  />
))
ToastAction.displayName = ToastPrimitives.Action.displayName

const ToastClose = React.forwardRef<
  React.ElementRef<typeof ToastPrimitives.Close>,
  React.ComponentPropsWithoutRef<typeof ToastPrimitives.Close>
>(({ className, ...props }, ref) => (
  <ToastPrimitives.Close
    ref={ref}
    className={cn(
      "absolute right-2 top-2 rounded-md p-1 text-foreground/50 opacity-70 transition-opacity hover:opacity-100 focus:opacity-100 focus:outline-none focus:ring-2 group-hover:opacity-100",
      "group-[.destructive]:text-destructive-foreground/70 group-[.destructive]:hover:text-destructive-foreground group-[.destructive]:focus:ring-destructive-foreground group-[.destructive]:focus:ring-offset-destructive", 
      "group-[.success]:text-green-400 group-[.success]:hover:text-green-600",
      "group-[.info]:text-blue-400 group-[.info]:hover:text-blue-600",
      "group-[.warning]:text-yellow-400 group-[.warning]:hover:text-yellow-600",
      className
    )}
    toast-close=""
    {...props}
  >
    <X className="h-4 w-4" />
  </ToastPrimitives.Close>
))
ToastClose.displayName = ToastPrimitives.Close.displayName

const ToastTitle = React.forwardRef<
  React.ElementRef<typeof ToastPrimitives.Title>,
  React.ComponentPropsWithoutRef<typeof ToastPrimitives.Title>
>(({ className, ...props }, ref) => (
  <ToastPrimitives.Title
    ref={ref}
    className={cn("text-sm font-semibold", className)}
    {...props}
  />
))
ToastTitle.displayName = ToastPrimitives.Title.displayName

const ToastDescription = React.forwardRef<
  React.ElementRef<typeof ToastPrimitives.Description>,
  React.ComponentPropsWithoutRef<typeof ToastPrimitives.Description>
>(({ className, ...props }, ref) => (
  <ToastPrimitives.Description
    ref={ref}
    className={cn("text-sm opacity-90", className)}
    {...props}
  />
))
ToastDescription.displayName = ToastPrimitives.Description.displayName

type ToastProps = React.ComponentPropsWithoutRef<typeof Toast>

type ToastActionElement = React.ReactElement<typeof ToastAction>

export {
  type ToastProps,
  type ToastActionElement,
  ToastProvider,
  ToastViewport,
  Toast,
  ToastTitle,
  ToastDescription,
  ToastClose,
  ToastAction,
}

