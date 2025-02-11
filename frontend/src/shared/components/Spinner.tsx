import React from 'react';
import { twMerge } from 'tailwind-merge';

export interface SpinnerProps extends React.SVGAttributes<SVGElement> {
  size?: 'sm' | 'md' | 'lg';
  label?: string;
}

export const Spinner = React.forwardRef<SVGSVGElement, SpinnerProps>(
  ({ className, size = 'md', label = 'Loading...', ...props }, ref) => {
    const sizes = {
      sm: 'h-4 w-4',
      md: 'h-8 w-8',
      lg: 'h-12 w-12',
    };

    return (
      <div
        role="status"
        className="inline-flex flex-col items-center justify-center"
      >
        <svg
          ref={ref}
          className={twMerge(
            'animate-spin',
            sizes[size],
            'text-muted-foreground',
            className
          )}
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          {...props}
        >
          <circle
            className="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            strokeWidth="4"
          />
          <path
            className="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
          />
        </svg>
        {label && (
          <span className="sr-only">
            {label}
          </span>
        )}
      </div>
    );
  }
);

Spinner.displayName = 'Spinner';

export const ProgressSpinner: React.FC<{ progress: number }> = ({ progress }) => {
  const circumference = 2 * Math.PI * 45; // r = 45
  const strokeDashoffset = circumference - (progress / 100) * circumference;

  return (
    <div className="relative inline-flex items-center justify-center">
      <svg className="h-24 w-24 transform -rotate-90" viewBox="0 0 100 100">
        {/* Background circle */}
        <circle
          className="text-muted stroke-current"
          strokeWidth="4"
          fill="transparent"
          r="45"
          cx="50"
          cy="50"
        />
        {/* Progress circle */}
        <circle
          className="text-primary stroke-current transition-all duration-300 ease-in-out"
          strokeWidth="4"
          strokeLinecap="round"
          fill="transparent"
          r="45"
          cx="50"
          cy="50"
          style={{
            strokeDasharray: circumference,
            strokeDashoffset: strokeDashoffset,
          }}
        />
      </svg>
      <div className="absolute text-lg font-semibold">
        {Math.round(progress)}%
      </div>
    </div>
  );
};
