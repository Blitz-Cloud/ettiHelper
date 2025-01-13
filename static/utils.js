async function getCodeToCopy(string) {
  toast = document.getElementById("toast-default");
  console.log(toast);
  toast.classList.remove("transition-opacity", "opacity-0", "hidden");
  try {
    const response = await fetch("/api/post/" + string);
    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }

    const data = await response.text();
    navigator.clipboard.writeText(data);
    setTimeout(() => {
      toast.classList.add("transition-opacity", "opacity-0", "hidden");
    }, 1500);
  } catch (error) {
    console.error(error.message);
  }
}

async function getCodeToCopyTipizat(string) {
  toast = document.getElementById("toast-default");
  console.log(toast);
  toast.classList.remove("transition-opacity", "opacity-0", "hidden");
  try {
    const response = await fetch("/api/tipizat/" + string);
    if (!response.ok) {
      throw new Error(`Response status: ${response.status}`);
    }

    const data = await response.text();
    navigator.clipboard.writeText(data);
    setTimeout(() => {
      toast.classList.add("transition-opacity", "opacity-0", "hidden");
    }, 1500);
  } catch (error) {
    console.error(error.message);
  }
}
