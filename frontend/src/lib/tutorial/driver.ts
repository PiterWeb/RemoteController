import { driver } from 'driver.js';
import type { DriveStep, Driver } from 'driver.js';
import 'driver.js/dist/driver.css';

// hazer un singletone para el tutorial
let tutorialDriver: Driver;

function StartTutorial(steps: DriveStep[], currentStep: number = 0) {
	if (tutorialDriver) {
		tutorialDriver.destroy();
		return tutorialDriver;
	}

	tutorialDriver = driver({
		animate: true,
		smoothScroll: true,
		stagePadding: 1,
		stageRadius: 1,
		onDestroyStarted: () => {
			if (!tutorialDriver.hasNextStep() || confirm('Are you sure?')) {
				tutorialDriver.destroy();
			}
		}
	});

	tutorialDriver.setSteps(steps);
	tutorialDriver.drive(currentStep);

	return tutorialDriver;
}

export default StartTutorial;
