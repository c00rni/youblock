Feature: Keep a list of domain block
    In order to stay focus on my work
    As a distracted guy
    I need to keep a list of website block

    Scenario: should block all domain in the list
        Given the list of block domain is:
            | youtube.com |
            | x.com |
        When I block the list of domain
        Then "youtube.com" can't be access anymore
        Then "x.com" can't be access anymore

    Scenario: should verfiy that domain list are block when the system change
        Given the list of block domain is:
            | youtube.com |
            | x.com |
        When the hosts file is modify
        Then verify that each domain is unreachable

    Scenario: should block new domains
        Given I want to block the domain "youtube.com" 
        And "youtube.com" is not already block
        When I block the list of domain
        Then the following domains can't be reached:
            | youtube.com |

