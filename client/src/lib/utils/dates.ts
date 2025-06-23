/**
 * Format a date string into a more readable format
 * @param dateString - ISO date string
 * @returns Formatted date string (e.g. "Jan 1, 2023")
 */
export function formatDate(dateString: string): string {
  if (!dateString) return "";

  const options: Intl.DateTimeFormatOptions = {
    year: "numeric",
    month: "short",
    day: "numeric",
  };

  return new Date(dateString).toLocaleDateString(undefined, options);
}

/**
 * Format a date string to include time
 * @param dateString - ISO date string
 * @returns Formatted date and time string (e.g. "Jan 1, 2023, 13:45")
 */
export function formatDateTime(dateString: string): string {
  if (!dateString) return "";

  const options: Intl.DateTimeFormatOptions = {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  };

  return new Date(dateString).toLocaleDateString(undefined, options);
}

/**
 * Get relative time from date (e.g. "2 days ago")
 * @param dateString - ISO date string
 * @returns Relative time string
 */
export function getRelativeTime(dateString: string): string {
  if (!dateString) return "";

  const now = new Date();
  const date = new Date(dateString);
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (diffInSeconds < 60) {
    return "just now";
  } else if (diffInSeconds < 3600) {
    const minutes = Math.floor(diffInSeconds / 60);
    return `${minutes} ${minutes === 1 ? "minute" : "minutes"} ago`;
  } else if (diffInSeconds < 86400) {
    const hours = Math.floor(diffInSeconds / 3600);
    return `${hours} ${hours === 1 ? "hour" : "hours"} ago`;
  } else if (diffInSeconds < 2592000) {
    const days = Math.floor(diffInSeconds / 86400);
    return `${days} ${days === 1 ? "day" : "days"} ago`;
  } else if (diffInSeconds < 31536000) {
    const months = Math.floor(diffInSeconds / 2592000);
    return `${months} ${months === 1 ? "month" : "months"} ago`;
  } else {
    const years = Math.floor(diffInSeconds / 31536000);
    return `${years} ${years === 1 ? "year" : "years"} ago`;
  }
}
