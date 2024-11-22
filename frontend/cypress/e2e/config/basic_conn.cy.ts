describe("Basic connection", () => {

    it("load", () => {
		cy.visit('http://localhost:34115/');
        cy.wait(1000)
        cy.log("hello")
    })

})