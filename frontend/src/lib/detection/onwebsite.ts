// This file is used to determine if the app is running on the website(client only) or not.
const onwebsite = import.meta.env?.VITE_ON_WEBSITE === 'true';

// This will be true if using the linux client when browser opens
const IS_RUNNING_EXTERNAL = window.location.port === "8080";

export {IS_RUNNING_EXTERNAL}

export default onwebsite;