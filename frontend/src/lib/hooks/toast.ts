import { writable } from "svelte/store";

interface Toast {
    show: boolean;
    message: string;
    type: 'info' | 'success' | 'error';
}

const toastWritable = writable<Toast>({
  show: false,
  message: "",
  type: "success",
});

export function showToast(message: string, type: Toast["type"]) {
  toastWritable.set({
    show: true,
    message,
    type,
  });

  setTimeout(() => {
    hideToast();
  }, 2000);
}

export function hideToast() {
  toastWritable.set({
    show: false,
    message: "",
    type: "success",
  });
}


export default toastWritable;