import { goto } from '$app/navigation';
import { driver } from 'driver.js';
import type { Driver, DriveStep } from 'driver.js';
import 'driver.js/dist/driver.css';
import { _ } from 'svelte-i18n';
import { get } from 'svelte/store';

// hazer un singletone para el tutorial
let tutorialDriver: Driver;
let currentStep = 0;

const TUTORIAL_DELAY = 750;

export function StartTutorial(selectedStep: number = 0) {
	if (tutorialDriver) {
		tutorialDriver.destroy();
	}

	tutorialDriver = driver({
		animate: true,
		smoothScroll: true,
		stagePadding: 1,
		stageRadius: 1,
		doneBtnText: get(_)('tutorial_done_text'),
		nextBtnText: get(_)('tutorial_next_text'),
		prevBtnText: get(_)('tutorial_prev_text')
	});

	const driverSteps: DriveStep[] = [
		{
			element: '#tutorial-config-btn',
			popover: {
				title: get(_)('tutorial_config_title'),
				description: get(_)('tutorial_config_description'),
				onNextClick: () => {
					goto('/mode/config');
					goNextTutorial();
				}
			}
		},
		{
			element: '#tutorial-language',
			popover: {
				title: get(_)('tutorial_language_title'),
				description: get(_)('tutorial_language_description'),
				onNextClick: () => {
					goNextTutorial();
				},
				onPrevClick: () => {
					goto('/');
					goPrevTutorial();
				}
			}
		},
		{
			element: '#tutorial-stun-card',
			popover: {
				title: get(_)('tutorial_stun_title'),
				description: get(_)('tutorial_stun_description'),
				onNextClick: () => {
					goto('/mode/config/advanced/stun');
					goNextTutorial();
				}
			}
		},
		{
			element: '#tutorial-group-server',
			popover: {
				title: get(_)('tutorial_group_server_title'),
				description: get(_)('tutorial_group_server_description'),
				onPrevClick: () => {
					goto('/mode/config');
					goPrevTutorial();
				},
				onNextClick: () => {
					goNextTutorial();
				}
			}
		},
		{
			element: '#tutorial-back-btn',
			popover: {
				title: get(_)('tutorial_go_back_title'),
				description: get(_)('tutorial_go_back_description'),
				onPrevClick: () => {
					goPrevTutorial();
				},
				onNextClick: () => {
					goto('/mode/config');
					goNextTutorial();
				}
			}
		},
		{
			element: '#tutorial-turn-card',
			popover: {
				title: get(_)('tutorial_turn_title'),
				description: get(_)('tutorial_turn_description'),
				onPrevClick: () => {
					goto('/mode/config/advanced/stun');
					goPrevTutorial();
				},
				onNextClick: () => {
					goNextTutorial();
				}
			}
		},
		{
			element: '#tutorial-back-btn',
			popover: {
				title: get(_)('tutorial_go_back_title'),
				description: get(_)('tutorial_go_back_description'),
				onPrevClick: () => {
					goPrevTutorial();
				},
				onNextClick: () => {
					goto('/');
					goNextTutorial();
				}
			}
		},
		{ 
			element: '#tutorial-play',
			popover: {
				title: get(_)('tutorial_play_title'),
				description: get(_)('tutorial_play_description'),
				onPrevClick: () => {
					goto("/mode/config")
					goPrevTutorial();
				},
			} 

		}
	];

	tutorialDriver.setSteps(driverSteps);
	tutorialDriver.drive(selectedStep);
}

function goNextTutorial(duration: number = TUTORIAL_DELAY) {
	setTimeout(() => {
		currentStep = currentStep + 1;
		tutorialDriver?.moveNext();
	}, duration);
}

function goPrevTutorial(duration: number = TUTORIAL_DELAY) {
	setTimeout(() => {
		currentStep = currentStep - 1;
		tutorialDriver?.movePrevious();
	}, duration);
}

_.subscribe(() => {
	if (tutorialDriver) {
		StartTutorial(currentStep);
	}
});
