describe('Tutorial flow', () => {
	it('Load main page ', () => {
		cy.visit('http://localhost:34115/');
	});

	it('tutorial', () => {

		cy.visit('http://localhost:34115/');
		cy.get('button.btn').click();
		cy.get('button.driver-popover-next-btn').click();
		cy.location('pathname').should('equal', '/mode/config');
		cy.wait(1000);
		cy.get('.grid > :nth-child(5) > .btn').click();
		cy.get(':nth-child(3) > .btn').click();
		cy.get(':nth-child(4) > .btn').click();
		cy.get('.grid > :nth-child(1) > .btn').click();
		cy.get('.grid > :nth-child(2) > .btn').click();
		cy.get('button.driver-popover-next-btn').click();
		cy.wait(1000);
		cy.get('button.driver-popover-next-btn').click();
		cy.location('pathname').should('equal', '/mode/config/advanced/stun');
		cy.wait(1000);
		cy.get('button.driver-popover-next-btn').click();
		cy.wait(1000);
		cy.get('button.driver-popover-next-btn').click();
		cy.wait(1000);
		cy.location('pathname').should('equal', '/mode/config');
		cy.get('button.driver-popover-next-btn').click();
		cy.wait(1000);
		cy.get('button.driver-popover-next-btn').click();
		cy.wait(1000);
		cy.get('button.driver-popover-next-btn').click();
		cy.location('pathname').should('equal', '/');
	});
});
