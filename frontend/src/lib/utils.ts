import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import { formatDistanceToNow } from "date-fns"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatDistance(meters: number): string {
  if (meters < 1000) {
    return `${Math.round(meters)}m`
  }
  return `${(meters / 1000).toFixed(1)}km`
}

export function formatTimeAgo(dateString: string): string {
  try {
    return formatDistanceToNow(new Date(dateString), { addSuffix: true })
  } catch (e) {
    return "just now"
  }
}

// Pulls the backend's { error: "..." } message out of an axios error so the
// UI can show the real reason (wrong password, duplicate email, rate limit,
// validation failure) instead of failing silently.
export function getErrorMessage(err: unknown, fallback = "Something went wrong. Please try again."): string {
  const message = (err as { response?: { data?: { error?: string } } })?.response?.data?.error
  return message || fallback
}
