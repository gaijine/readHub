package service

import (
	"log"

	openlibrary "readHub/internal/client/openlibrary"
	"readHub/internal/domain"
	"readHub/internal/repository"
)

type BookService interface {
	SearchBooks(query string) ([]domain.SearchBook, error)
	AddBook(userID int64, bookInfo domain.SearchBook) error
	GetUserBooks(userID int64) ([]domain.Book, error)
	UpdateStatus(userID, bookID int64, status domain.BookStatus) error
	UpdateProgress(userID, bookID int64, page int) error
	DeleteBook(userID, bookID int64) error
	GetBookDetails(openLibraryID string) (domain.BookDetails, error)
	// AddBookByOpenLibraryID(telegramID int64, openLibraryID string) error
	// GetSearchBookByID(openLibraryID string)(domain.SearchBook, error)
	GetUserByTelegramID(telegramID int64) (domain.User, error)
	CreateUser(user domain.User) error
}

type bookService struct {
	bookRepo repository.BookRepository // здесь лежит объект который умеет работать с книгами
	userRepo repository.UserRepository
	openLib  openlibrary.OpenLibraryClient
}

// функция должна вернуть BookService.
// (*bookService реализует все методы BookService.)
// поэтому Go разрешает вернуть *bookService (&bookService) как BookService.
func NewBookService(bookRepo repository.BookRepository, userRepo repository.UserRepository, openLib openlibrary.OpenLibraryClient) BookService {
	return &bookService{
		bookRepo: bookRepo,
		userRepo: userRepo,
		openLib:  openLib,
	}
}

// это одно и то же, получается переменная сервис с типом BookService, как сказать
// получается может принять любой объект который реализует методы и удовлетворит интерфейс BookService
// и вто же время если посмотреть то данная структура в которой хранятся другие интерфейсы реализует методы
//{
/*  var service BookService = bookService{
	  bookRepo: bookRepo,
	  userRepo: userRepo,
	  openLib: openLib,
   }
    return service*/
//}

func (b *bookService) SearchBooks(query string) ([]domain.SearchBook, error) {
	books, err := b.openLib.SearchBooks(query)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (b *bookService) AddBook(userID int64, bookInfo domain.SearchBook) error {
	// bookInfo, err := b.openLib.GetByOpenLibraryID(openLibraryID)
	// if err != nil {
	// 	return err
	// }

	var author string
	if len(bookInfo.Author) == 0 {
		author = "Unknown"
	} else {
		author = bookInfo.Author[0]
	}

	book := domain.Book{
		UserID:        userID,
		OpenLibraryID: bookInfo.OpenLibraryID,
		Title:         bookInfo.Title,
		Author:        author,
		TotalPages:    0,
		CurrentPage:   0,
		Status:        domain.StatusWantToRead,
		CoverURL:      bookInfo.CoverURL,
	}

	// return b.bookRepo.Create(book) так тоже можно записать, мол если метод вернет nil, то функция также вернет nil, с err также
	err := b.bookRepo.Create(book)
	if err != nil {
		return err
	}
	log.Println("book added successed")
	return nil
}

func (b *bookService) GetUserBooks(userID int64) ([]domain.Book, error) {
	// return b.bookRepo.GetByUserID(userID) также можно сделать, так как возвращают одно и тоже
	books, err := b.bookRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (b *bookService) UpdateStatus(userID, bookID int64, status domain.BookStatus) error {
	book, err := b.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}

	if book.UserID != userID {
		return ErrBookNotOwned
	}

	err = b.bookRepo.UpdateStatus(bookID, status)
	if err != nil {
		return err
	}
	return nil
}

func (b *bookService) UpdateProgress(userID, bookID int64, page int) error {
	book, err := b.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}

	if book.UserID != userID {
		return ErrBookNotOwned
	}

	if page < 0 {
		return ErrNegativePage
	}

	if page > book.TotalPages {
		return ErrPageExceedsBookLength
	}

	err = b.bookRepo.UpdateCurrentPage(bookID, page)
	if err != nil {
		return err
	}

	if page == book.TotalPages {
		err = b.bookRepo.UpdateStatus(bookID, domain.StatusCompleted)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *bookService) DeleteBook(userID, bookID int64) error {
	book, err := b.bookRepo.GetByID(bookID)
	if err != nil {
		return err
	}

	if userID != book.UserID {
		return ErrBookNotOwned
	}

	err = b.bookRepo.Delete(bookID)
	if err != nil {
		return err
	}
	return nil
}

func (b *bookService) GetBookDetails(openLibraryID string) (domain.BookDetails, error) {
	book, err := b.openLib.GetBookDetails(openLibraryID)
	if err != nil {
		return domain.BookDetails{}, err
	}
	return book, nil
}

// func (b *bookService) AddBookByOpenLibraryID(teledramID int64, openLibraryID string) error {

// }

// func (b *bookService) GetSearchBookByID (openLibraryID string)(domain.SearchBook, error){

// }

func (b *bookService) GetUserByTelegramID(telegramID int64) (domain.User, error) {
	user, err := b.userRepo.GetByTelegramID(telegramID)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (b *bookService) CreateUser(user domain.User) error {
	err := b.userRepo.Create(user)
	if err != nil {
		return err
	}
	return nil
}
