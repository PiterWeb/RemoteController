// This file is used to determine if the app is running on the website(client only) or not.
const onwebsite = import.meta.env?.VITE_ON_WEBSITE === 'true';

export default onwebsite;