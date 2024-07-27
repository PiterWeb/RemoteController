import { driver } from 'driver.js';
import 'driver.js/dist/driver.css';

// hazer un singletone para el tutorial

const TutorialDriver = driver({
	animate: true,
    smoothScroll: true,
    stagePadding: 1,
	stageRadius: 1,
	onDestroyStarted: () => {
		if (!TutorialDriver.hasNextStep() || confirm('Are you sure?')) {
			TutorialDriver.destroy();
		}
	}
});

export default TutorialDriver;
