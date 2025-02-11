import React from 'react';
import { twMerge } from 'tailwind-merge';

export interface SwitchProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'type'> {
  label?: string;
  description?: string;
  error?: string;
}

export const Switch = React.forwardRef<HTMLInputElement, SwitchProps>(
  (
    {
      className,
      label,
      description,
      error,
      disabled,
      id,
      ...props
    },
    ref
  ) => {
    const switchId = id || label?.toLowerCase().replace(/\s+/g, '-');

    return (
      <div className={twMerge('relative flex items-start', className)}>
        <div className="flex h-6 items-center">
          <input
            type="checkbox"
            ref={ref}
            id={switchId}
            disabled={disabled}
            aria-invalid={!!error}
            aria-describedby={
              error
                ? `${switchId}-error`
                : description
                ? `${switchId}-description`
                : undefined
            }
            className={twMerge(
              'h-4 w-7 rounded-full appearance-none bg-muted transition-colors',
              'checked:bg-primary peer relative cursor-pointer',
              'disabled:cursor-not-allowed disabled:opacity-50',
              'after:content-[""] after:block after:h-3 after:w-3',
              'after:rounded-full after:bg-background after:transition-transform',
              'after:absolute after:top-0.5 after:left-0.5',
              'checked:after:translate-x-3',
              'focus-visible:outline-none focus-visible:ring-2',
              'focus-visible:ring-ring focus-visible:ring-offset-2',
              error && 'border-destructive focus-visible:ring-destructive'
            )}
            {...props}
          />
        </div>
        {(label || description) && (
          <div className="ml-3">
            {label && (
              <label
                htmlFor={switchId}
                className={twMerge(
                  'text-sm font-medium leading-6',
                  disabled && 'cursor-not-allowed opacity-50'
                )}
              >
                {label}
              </label>
            )}
            {description && (
              <p
                id={`${switchId}-description`}
                className="text-sm text-muted-foreground"
              >
                {description}
              </p>
            )}
            {error && (
              <p
                id={`${switchId}-error`}
                className="text-sm text-destructive"
                role="alert"
              >
                {error}
              </p>
            )}
          </div>
        )}
      </div>
    );
  }
);

Switch.displayName = 'Switch';
