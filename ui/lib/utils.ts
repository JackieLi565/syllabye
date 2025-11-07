import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";
import CRC32 from "crc-32";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function schoolYearFormatter(year: number | null | undefined) {
  if (!year) {
    return "No year provided";
  }

  switch (year) {
    case 1:
      return "1st Year";
    case 2:
      return "2nd Year";
    case 3:
      return "3rd Year";
    default:
      return `${year}th Year`;
  }
}

export const formatFileSize = (bytes: number) => {
  if (bytes < 1024) return bytes + " B";
  else if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + " KB";
  else return (bytes / (1024 * 1024)).toFixed(1) + " MB";
};

export const calculateChecksum = async (file: File): Promise<string> => {
  const arrayBuffer = await file.arrayBuffer();
  const uint8Array = new Uint8Array(arrayBuffer);
  const checksum = CRC32.buf(uint8Array);
  return (checksum >>> 0).toString(16);
};

export const getFirstQueryParam = <T>(param: T | T[]) => {
  let value;
  if (Array.isArray(param)) {
    value = param[0];
  } else {
    value = param;
  }

  return value !== null ? value : undefined;
};
