import { ButtonHTMLAttributes } from 'react';
import { cn } from '@/lib/utils/cn';

export function Button({ className, ...props }: ButtonHTMLAttributes<HTMLButtonElement>) {
  return <button className={cn('px-3 py-2 bg-blue-600 text-white', className)} {...props} />;
}
