/******************************************************************************
 * =============================================================================
 * BAOBAB ENTERPRISE PLATFORM
 * =============================================================================
 *
 * File:
 *      docs/assets/js/extra.js
 *
 * Purpose:
 *      Progressive enhancements for the Baobab Documentation Portal.
 *
 * Description:
 *      This file provides lightweight client-side enhancements that complement
 *      Material for MkDocs without overriding its native functionality.
 *
 *      The documentation site must remain fully functional even if JavaScript
 *      is disabled. Every enhancement in this file is optional and improves
 *      usability, accessibility, and user experience.
 *
 * Responsibilities:
 *
 *      • Application bootstrap
 *      • Event registration
 *      • Smooth scrolling
 *      • External link enhancements
 *      • Code copy feedback
 *      • Scroll-to-top button
 *      • Reveal animations
 *      • Hero enhancements
 *      • Utility functions
 *      • Accessibility helpers
 *
 * Design Principles:
 *
 *      • Progressive Enhancement
 *      • Accessibility First
 *      • Performance Focused
 *      • Minimal DOM Manipulation
 *      • No Framework Dependencies
 *      • No Theme Overrides
 *
 * -----------------------------------------------------------------------------
 * Future Enhancements
 * -----------------------------------------------------------------------------
 *
 * The following features are intentionally planned for future releases and may
 * later be extracted into dedicated JavaScript modules.
 *
 * Navigation
 * ----------
 * • Active section highlighting
 * • Sticky navigation enhancements
 * • Breadcrumb improvements
 * • Keyboard navigation
 *
 * Search
 * ------
 * • Search suggestions
 * • Recent searches
 * • Search result highlighting
 * • Search analytics
 *
 * Documentation
 * -------------
 * • Mermaid diagram helpers
 * • PlantUML enhancements
 * • Interactive API examples
 * • OpenAPI playground
 * • UML diagram zoom controls
 *
 * Clipboard
 * ---------
 * • Copy page URL
 * • Copy section anchors
 * • Copy API endpoint
 * • Copy command snippets
 *
 * Accessibility
 * -------------
 * • Skip navigation improvements
 * • Focus management
 * • Keyboard shortcuts
 * • High contrast support
 *
 * Theme
 * -----
 * • Theme preference persistence
 * • Custom colour palettes
 * • Reading mode
 * • Presentation mode
 *
 * Analytics
 * ---------
 * • Documentation usage metrics
 * • Page performance metrics
 * • Broken link detection
 * • Client diagnostics
 *
 * Utilities
 * ---------
 * • Lazy loading helpers
 * • Image lightbox
 * • Notification system
 * • Toast messages
 * • Scroll spy
 * • Page loader
 * • Service Worker support
 * • Offline documentation
 *
 ******************************************************************************/

"use strict";

/*==============================================================================
    Bootstrap
==============================================================================*/

document.addEventListener("DOMContentLoaded", initialize);


/*==============================================================================
    Application Initialisation
==============================================================================*/

function initialize() {

    enableSmoothScrolling();

    configureExternalLinks();

    enhanceClipboard();

    createScrollToTopButton();

    initialiseRevealAnimations();

    initialiseHero();

    console.info("Baobab Documentation Portal initialised.");

}


/*==============================================================================
    Smooth Scrolling
==============================================================================*/

function enableSmoothScrolling() {

    document.documentElement.style.scrollBehavior = "smooth";

}


/*==============================================================================
    External Links
==============================================================================*/

function configureExternalLinks() {

    document
        .querySelectorAll('a[href^="http"]')
        .forEach(link => {

            if (link.hostname !== window.location.hostname) {

                link.target = "_blank";

                link.rel = "noopener noreferrer";

            }

        });

}


/*==============================================================================
    Clipboard Feedback
==============================================================================*/

function enhanceClipboard() {

    document
        .querySelectorAll(".md-clipboard")
        .forEach(button => {

            button.addEventListener("click", () => {

                button.classList.add("copied");

                setTimeout(() => {

                    button.classList.remove("copied");

                }, 1200);

            });

        });

}


/*==============================================================================
    Scroll To Top
==============================================================================*/

function createScrollToTopButton() {

    const button = document.createElement("button");

    button.className = "scroll-top";

    button.type = "button";

    button.setAttribute("aria-label", "Scroll to top");

    button.innerHTML = "↑";

    document.body.appendChild(button);

    window.addEventListener("scroll", () => {

        if (window.scrollY > 500) {

            button.classList.add("visible");

        }

        else {

            button.classList.remove("visible");

        }

    });

    button.addEventListener("click", () => {

        window.scrollTo({

            top: 0,

            behavior: "smooth"

        });

    });

}


/*==============================================================================
    Reveal Animations
==============================================================================*/

function initialiseRevealAnimations() {

    if (!("IntersectionObserver" in window)) {

        return;

    }

    const observer = new IntersectionObserver(

        entries => {

            entries.forEach(entry => {

                if (entry.isIntersecting) {

                    entry.target.classList.add("fade-in");

                    observer.unobserve(entry.target);

                }

            });

        },

        {

            threshold: 0.15

        }

    );

    document
        .querySelectorAll(".grid.cards li")
        .forEach(card => observer.observe(card));

}


/*==============================================================================
    Hero
==============================================================================*/

function initialiseHero() {

    const hero = document.querySelector(".hero");

    if (!hero) {

        return;

    }

    hero.classList.add("slide-up");

}


/*==============================================================================
    Utility
==============================================================================*/

function debounce(callback, delay = 250) {

    let timer;

    return (...args) => {

        clearTimeout(timer);

        timer = setTimeout(() => {

            callback(...args);

        }, delay);

    };

}


/*==============================================================================
    Window Resize
==============================================================================*/

window.addEventListener(

    "resize",

    debounce(() => {

        // Reserved for future responsive enhancements.

    })

);


/*==============================================================================
    Error Handling
==============================================================================*/

window.addEventListener("error", event => {

    console.error(

        "Baobab documentation error:",

        event.message

    );

});


/*==============================================================================
    End of File
==============================================================================*/
