import { type ClassValue, clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function schoolYearFormatter(year: number | null | undefined) {
  if (!year) {
    return "No year provided"
  }

  switch (year) {
    case 1:
      return "1st Year"
    case 2:
      return "2nd Year"
    case 3:
      return "3rd Year"
    default:
      return `${year}th Year`
  }
}