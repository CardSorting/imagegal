import React from 'react';
import { twMerge } from 'tailwind-merge';

export interface SliderProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'type'> {
  label?: string;
  error?: string;
  helperText?: string;
  showValue?: boolean;
  valuePrefix?: string;
  valueSuffix?: string;
}

export const Slider = React.forwardRef<HTMLInputElement, SliderProps>(
  (
    {
      className,
      label,
      error,
      helperText,
      showValue = true,
      valuePrefix = '',
      valueSuffix = '',
      value,
      min,
      max,
      step,
      disabled,
      id,
      ...props
    },
    ref
  ) => {
    const inputId = id || label?.toLowerCase().replace(/\s+/g, '-');

    return (
      <div className="w-full">
        <div className="flex items-center justify-between mb-1.5">
          {label && (
            <label
              htmlFor={inputId}
              className="block text-sm font-medium text-foreground"
            >
              {label}
            </label>
          )}
          {showValue && (
            <span className="text-sm text-muted-foreground">
              {valuePrefix}{value}{valueSuffix}
            </span>
          )}
        </div>
        <div className="relative">
          <input
            ref={ref}
            type="range"
            id={inputId}
            value={value}
            min={min}
            max={max}
            step={step}
            disabled={disabled}
            className={twMerge(
              'w-full h-2 bg-input rounded-lg appearance-none cursor-pointer',
              'dark:bg-gray-800',
              '[&::-webkit-slider-thumb]:appearance-none',
              '[&::-webkit-slider-thumb]:w-4',
              '[&::-webkit-slider-thumb]:h-4',
              '[&::-webkit-slider-thumb]:rounded-full',
              '[&::-webkit-slider-thumb]:bg-primary hover:[&::-webkit-slider-thumb]:bg-primary/90',
              '[&::-webkit-slider-thumb]:cursor-pointer',
              '[&::-moz-range-thumb]:w-4',
              '[&::-moz-range-thumb]:h-4',
              '[&::-moz-range-thumb]:rounded-full',
              '[&::-moz-range-thumb]:bg-primary hover:[&::-moz-range-thumb]:bg-primary/90',
              '[&::-moz-range-thumb]:cursor-pointer',
              '[&::-moz-range-thumb]:border-0',
              disabled && 'opacity-50 cursor-not-allowed',
              error && 'bg-destructive/20',
              className
            )}
            {...props}
          />
        </div>
        {(error || helperText) && (
          <div className="mt-1.5">
            {error && (
              <p
                className="text-sm text-destructive"
                id={`${inputId}-error`}
                role="alert"
              >
                {error}
              </p>
            )}
            {!error && helperText && (
              <p className="text-sm text-muted-foreground">{helperText}</p>
            )}
          </div>
        )}
      </div>
    );
  }
);

Slider.displayName = 'Slider';
