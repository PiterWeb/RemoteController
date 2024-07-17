import { driver } from 'driver.js';
import 'driver.js/dist/driver.css';

const TutorialDriver = driver({
	animate: true,
    smoothScroll: true,
    stagePadding: 1,
	onDestroyStarted: () => {
		if (!TutorialDriver.hasNextStep() || confirm('Are you sure?')) {
			TutorialDriver.destroy();
		}
	}
});

export default TutorialDriver;
