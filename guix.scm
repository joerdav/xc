(define-module (guix)
  #:use-module (guix packages)
  #:use-module (guix gexp)
  #:use-module (guix git-download)
  #:use-module (gnu packages task-runners)
  #:use-module (guix utils))

(define-public xc-dev
  (package
   (inherit xc)
   (version "0.dev")
   (source
    (local-file
     "." "xc-checkout"
     #:recursive? #t
     #:select? (lambda* (#:rest _)
                 (or (git-predicate (current-source-directory))
                     (const #t)))))))

xc-dev
