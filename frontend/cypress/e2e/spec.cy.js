describe('Todo App Test', () => {
  it('Visits todo app', () => {
    cy.visit('http://127.0.0.1:5173')

    cy.get('input[name="task"]').type('play').should('have.value','play')
    cy.get('form').submit()
    cy.get('button[name="play-toggle"]').click()
    cy.get('button[name="play-delete"]').click()

    cy.get('input[name="task"]').type('run').should('have.value','run')
    cy.get('form').submit()
    cy.get('button[name="run-toggle"]').click()
    cy.get('button[name="run-delete"]').click()

    cy.get('input[name="task"]').type('work').should('have.value','work')
    cy.get('form').submit()
    cy.get('button[name="work-toggle"]').click()
    cy.get('button[name="work-delete"]').click()
  })
})