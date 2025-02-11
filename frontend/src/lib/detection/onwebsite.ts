// This file is used to determine if the app is running on the website(client only) or not.
const onwebsite = import.meta.env?.VITE_ON_WEBSITE === 'true';

const IS_RUNNING_EXTERNAL = window.location.port === "8081";

export {IS_RUNNING_EXTERNAL}

export default onwebsite;